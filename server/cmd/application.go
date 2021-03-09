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
	"fmt"
	"net/http"

	cloudstorage "cloud.google.com/go/storage"
	"github.com/batect/abacus/server/api"
	"github.com/batect/abacus/server/storage"
	"github.com/batect/service-observability/graceful"
	"github.com/batect/service-observability/middleware"
	"github.com/batect/service-observability/startup"
	"github.com/batect/service-observability/tracing"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

func main() {
	flush, err := startup.InitialiseObservability(getServiceName(), getVersion(), getProjectID())

	if err != nil {
		logrus.WithError(err).Fatal("Could not initialise observability.")
	}

	defer flush()

	srv := createServer(getPort())
	graceful.RunServerWithGracefulShutdown(srv)
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
			otelhttp.WithSpanNameFormatter(tracing.NameHTTPRequestSpan),
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
