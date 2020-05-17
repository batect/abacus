// Copyright 2019-2020 Charles Korn.
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
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/batect/abacus/server/api"
	"github.com/batect/abacus/server/middleware"
	"github.com/batect/abacus/server/observability"
	"github.com/batect/abacus/server/storage"
	stackdriver "github.com/charleskorn/logrus-stackdriver-formatter"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/plugin/othttp"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/api/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	initLogging()
	initTracing()

	srv := createServer(getPort())
	runServer(srv)
}

func initLogging() {
	logrus.SetFormatter(stackdriver.NewFormatter(
		stackdriver.WithService(getEnvOrDefault("K_SERVICE", "abacus")),
		stackdriver.WithVersion(getEnvOrDefault("K_REVISION", "local")),
	))
}

func getEnvOrDefault(name string, fallback string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}

	return fallback
}

func initTracing() {
	exporter, err := texporter.NewExporter(texporter.WithOnError(func(err error) {
		logrus.WithError(err).Warn("Trace exporter reported error.")
	}))

	if err != nil {
		logrus.WithError(err).Fatal("Could not create trace exporter.")
	}

	traceProvider, err := sdktrace.NewProvider(sdktrace.WithSyncer(exporter))

	if err != nil {
		logrus.WithError(err).Fatal("Could not create trace provider.")
	}

	global.SetTraceProvider(traceProvider)
	global.SetPropagators(&observability.GCPPropagator{})
}

func createServer(port string) *http.Server {
	mux := http.NewServeMux()

	// TODO: decorate these with othttp.WithRouteTag()
	mux.HandleFunc("/ping", api.Ping)
	mux.Handle("/v1/sessions", othttp.WithRouteTag("/v1/sessions", createIngestHandler()))

	wrappedMux := middleware.TraceIDExtractionMiddleware(
		middleware.LoggerMiddleware(logrus.StandardLogger(), getProjectID(), mux),
	)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: othttp.NewHandler(
			wrappedMux,
			"Incoming API call",
			othttp.WithMessageEvents(othttp.ReadEvents, othttp.WriteEvents),
			othttp.WithSpanNameFormatter(observability.NameHTTPRequestSpan),
		),
	}

	return srv
}

func createIngestHandler() http.Handler {
	store, err := storage.NewBigQuerySessionStore(getProjectID(), getDatasetID(), getSessionsTableID(), getCredentialsFilePath())

	if err != nil {
		logrus.WithError(err).Fatal("Could not create session store.")
	}

	return api.NewIngestHandler(store)
}

func runServer(srv *http.Server) {
	connectionDrainingFinished := shutdownOnInterrupt(srv)

	logrus.WithField("address", srv.Addr).Info("Server starting.")

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logrus.WithError(err).Fatal("Could not start HTTP server.")
	}

	<-connectionDrainingFinished

	logrus.Info("Server shut down.")
}

func shutdownOnInterrupt(srv *http.Server) chan struct{} {
	connectionDrainingFinished := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		logrus.Info("Interrupt received, draining connections...")

		if err := srv.Shutdown(context.Background()); err != nil {
			logrus.WithError(err).Error("Shutting down HTTP server failed.")
		}

		close(connectionDrainingFinished)
	}()

	return connectionDrainingFinished
}
