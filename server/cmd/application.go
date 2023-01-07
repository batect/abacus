// Copyright 2019-2023 Charles Korn.
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
	"time"

	cloudstorage "cloud.google.com/go/storage"
	"github.com/batect/abacus/server/api"
	"github.com/batect/abacus/server/storage"
	"github.com/batect/services-common/graceful"
	"github.com/batect/services-common/middleware"
	"github.com/batect/services-common/startup"
	"github.com/batect/services-common/tracing"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

func main() {
	config, err := getConfig()

	if err != nil {
		logrus.WithError(err).Error("Could not load application configuration.")
		os.Exit(1)
	}

	flush, err := startup.InitialiseObservability(config.ServiceName, config.ServiceVersion, config.ProjectID, config.HoneycombAPIKey)

	if err != nil {
		logrus.WithError(err).Error("Could not initialise observability tooling.")
		os.Exit(1)
	}

	defer flush()

	runServer(config)
}

func runServer(config *serviceConfig) {
	srv, err := createServer(config)

	if err != nil {
		logrus.WithError(err).Error("Could not create server.")
		os.Exit(1)
	}

	if err := graceful.RunServerWithGracefulShutdown(srv); err != nil {
		logrus.WithError(err).Error("Could not run server.")
		os.Exit(1)
	}
}

func createServer(config *serviceConfig) (*http.Server, error) {
	mux := http.NewServeMux()
	mux.Handle("/", otelhttp.WithRouteTag("/", http.HandlerFunc(api.Home)))
	mux.Handle("/ping", otelhttp.WithRouteTag("/ping", http.HandlerFunc(api.Ping)))

	ingestHandler, err := createIngestHandler(config)

	if err != nil {
		return nil, fmt.Errorf("could not create ingest endpoint handler: %w", err)
	}

	mux.Handle("/v1/sessions", otelhttp.WithRouteTag("/v1/sessions", ingestHandler))

	securityHeaders := secure.New(secure.Options{
		FrameDeny:             true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'none'; frame-ancestors 'none'",
		ReferrerPolicy:        "no-referrer",
	})

	wrappedMux := middleware.TraceIDExtractionMiddleware(
		middleware.LoggerMiddleware(
			logrus.StandardLogger(),
			config.ProjectID,
			securityHeaders.Handler(mux),
		),
	)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", config.Port),
		Handler: otelhttp.NewHandler(
			wrappedMux,
			"Abacus",
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
			otelhttp.WithSpanNameFormatter(tracing.NameHTTPRequestSpan),
		),
		ReadHeaderTimeout: 10 * time.Second,
	}

	return srv, nil
}

func createIngestHandler(config *serviceConfig) (http.Handler, error) {
	scopesOption := option.WithScopes(cloudstorage.ScopeReadWrite)
	credsOption := option.WithCredentialsFile(getCredentialsFilePath())
	tracingClientOption, err := withTracingClient(scopesOption, credsOption)

	if err != nil {
		return nil, fmt.Errorf("could not create tracing client: %w", err)
	}

	bucketName := fmt.Sprintf("%v-sessions", config.ProjectID)
	store, err := storage.NewCloudStorageSessionStore(bucketName, tracingClientOption)

	if err != nil {
		return nil, fmt.Errorf("could not create session store: %w", err)
	}

	handler, err := api.NewIngestHandler(store)

	if err != nil {
		return nil, fmt.Errorf("could not instantiate ingest API handler: %w", err)
	}

	return handler, nil
}

func withTracingClient(opts ...option.ClientOption) (option.ClientOption, error) {
	// We have to do this because setting http.DefaultTransport to a non-default implementation causes something deep in the bowels of the
	// Google Cloud SDK to ignore it and create a fresh transport with many of the settings copied across from DefaultTransport.
	// Being explicit about the client forces the SDK to use the transport.
	trans, err := htransport.NewTransport(context.Background(), http.DefaultTransport, opts...)

	if err != nil {
		return nil, fmt.Errorf("could not create transport: %w", err)
	}

	httpClient := http.Client{
		Transport: trans,
	}

	return option.WithHTTPClient(&httpClient), nil
}
