// Copyright 2019-2021 Charles Korn.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// and the Commons Clause License Condition v1.0 (the "Condition");
// you may not use this file except in compliance with both the License and Condition.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// You may obtain a copy of the Condition at
//
//     https://commonsclause.com/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License and the Condition is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See both the License and the Condition for the specific language governing permissions and
// limitations under the License and the Condition.

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/profiler"
	cloudstorage "cloud.google.com/go/storage"
	"github.com/batect/abacus/server/api"
	"github.com/batect/abacus/server/middleware"
	"github.com/batect/abacus/server/observability"
	"github.com/batect/abacus/server/storage"
	stackdriver "github.com/charleskorn/logrus-stackdriver-formatter"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/api/option"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	htransport "google.golang.org/api/transport/http"
)

func main() {
	initLogging()
	initProfiling()

	flush := initTracing()

	defer func() {
		logrus.Info("Flushing remaining traces...")
		flush()
		logrus.Info("Flushing complete.")
	}()

	srv := createServer(getPort())
	runServer(srv)
}

func initLogging() {
	logrus.SetFormatter(stackdriver.NewFormatter(
		stackdriver.WithService(getServiceName()),
		stackdriver.WithVersion(getVersion()),
	))
}

func getServiceName() string {
	return getEnvOrDefault("K_SERVICE", "abacus")
}

func getVersion() string {
	return getEnvOrDefault("K_REVISION", "local")
}

func getEnvOrDefault(name string, fallback string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}

	return fallback
}

func initProfiling() {
	err := profiler.Start(profiler.Config{
		Service:        getServiceName(),
		ServiceVersion: getVersion(),
		ProjectID:      getProjectID(),
		MutexProfiling: true,
	})

	if err != nil {
		logrus.WithError(err).Fatal("Could not create profiler.")
	}
}

func initTracing() func() {
	_, flush, err := texporter.InstallNewPipeline(
		[]texporter.Option{
			texporter.WithProjectID(getProjectID()),
			texporter.WithOnError(func(err error) {
				logrus.WithError(err).Warn("Trace exporter reported error.")
			}),
		},
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)

	if err != nil {
		logrus.WithError(err).Fatal("Could not install tracing pipeline.")
	}

	w3Propagator := propagation.TraceContext{}
	gcpPropagator := observability.GCPPropagator{}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(w3Propagator, gcpPropagator))

	http.DefaultTransport = otelhttp.NewTransport(
		http.DefaultTransport,
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		otelhttp.WithSpanNameFormatter(observability.NameHTTPRequestSpan),
	)

	return flush
}

func createServer(port string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/ping", otelhttp.WithRouteTag("/ping", http.HandlerFunc(api.Ping)))
	mux.Handle("/v1/sessions", otelhttp.WithRouteTag("/v1/sessions", createIngestHandler()))

	securityHeaders := secure.New(secure.Options{
		FrameDeny:             true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'none'; frame-ancestors 'none'",
		ReferrerPolicy:        "no-referrer",
	})

	wrappedMux := middleware.TraceIDExtractionMiddleware(
		middleware.LoggerMiddleware(
			logrus.StandardLogger(),
			getProjectID(),
			securityHeaders.Handler(mux),
		),
	)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: otelhttp.NewHandler(
			wrappedMux,
			"Incoming API call",
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
			otelhttp.WithSpanNameFormatter(observability.NameHTTPRequestSpan),
		),
	}

	return srv
}

func createIngestHandler() http.Handler {
	scopesOption := option.WithScopes(cloudstorage.ScopeReadWrite)
	credsOption := option.WithCredentialsFile(getCredentialsFilePath())
	tracingClientOption := withTracingClient(scopesOption, credsOption)
	bucketName := fmt.Sprintf("%v-sessions", getProjectID())
	store, err := storage.NewCloudStorageSessionStore(bucketName, tracingClientOption)

	if err != nil {
		logrus.WithError(err).Fatal("Could not create session store.")
	}

	handler, err := api.NewIngestHandler(store)

	if err != nil {
		logrus.WithError(err).Fatal("Could not create ingest API handler.")
	}

	return handler
}

func withTracingClient(opts ...option.ClientOption) option.ClientOption {
	// We have to do this because setting http.DefaultTransport to a non-default implementation causes something deep in the bowels of the
	// Google Cloud SDK to ignore it and create a fresh transport with many of the settings copied across from DefaultTransport.
	// Being explicit about the client forces the SDK to use the transport.
	trans, err := htransport.NewTransport(context.Background(), http.DefaultTransport, opts...)

	if err != nil {
		logrus.WithError(err).Fatal("could not create transport")
	}

	httpClient := http.Client{
		Transport: trans,
	}

	return option.WithHTTPClient(&httpClient)
}

func runServer(srv *http.Server) {
	connectionDrainingFinished := shutdownOnInterrupt(srv)

	logrus.WithField("address", srv.Addr).Info("Server starting.")

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logrus.WithError(err).Fatal("Could not start HTTP server.")
	}

	<-connectionDrainingFinished

	logrus.Info("Server gracefully stopped.")
}

func shutdownOnInterrupt(srv *http.Server) chan struct{} {
	connectionDrainingFinished := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		logrus.Info("Interrupt received, draining connections...")

		if err := srv.Shutdown(context.Background()); err != nil {
			logrus.WithError(err).Error("Shutting down HTTP server failed.")
		}

		close(connectionDrainingFinished)
	}()

	return connectionDrainingFinished
}
