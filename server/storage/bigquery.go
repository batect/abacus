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

package storage

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/batect/abacus/server/observability"
	"go.opentelemetry.io/otel/plugin/othttp"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
)

type bigQuerySessionStore struct {
	inserter *bigquery.Inserter
}

func NewBigQuerySessionStore(projectID string, datasetID string, tableID string, credsFile string) (SessionStore, error) {
	ctx := context.Background()

	// We have to do this because specifying option.WithHTTPClient in the call to bigquery.NewClient overrides all other options -
	// so instead we create the transport ourselves and then wrap that in the OpenTelemetry transport.
	// Would be good to investigate just setting http.DefaultTransport to be the OpenTelemetry transport so all HTTP calls get telemetry.
	scopesOption := option.WithScopes("https://www.googleapis.com/auth/bigquery.insertdata")
	credsOption := option.WithCredentialsFile(credsFile)
	trans, err := htransport.NewTransport(ctx, http.DefaultTransport, scopesOption, credsOption)

	if err != nil {
		return nil, fmt.Errorf("could not create transport: %w", err)
	}

	httpClient := http.Client{
		Transport: othttp.NewTransport(
			trans,
			othttp.WithMessageEvents(othttp.ReadEvents, othttp.WriteEvents),
			othttp.WithSpanNameFormatter(observability.NameHTTPRequestSpan),
		),
	}

	client, err := bigquery.NewClient(ctx, projectID, option.WithHTTPClient(&httpClient))

	if err != nil {
		return nil, err
	}

	dataset := client.Dataset(datasetID)
	table := dataset.Table(tableID)
	inserter := table.Inserter()

	return &bigQuerySessionStore{inserter}, nil
}

func (b *bigQuerySessionStore) Store(ctx context.Context, session *Session) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if err := b.inserter.Put(ctx, session); err != nil {
		if e, ok := err.(bigquery.PutMultiError); ok && len(e) == 1 {
			return fmt.Errorf("could not store session due to Put error: %w", &e[0])
		}

		return fmt.Errorf("could not store session: %w", err)
	}

	return nil
}

// Ensure that Session implements ValueSaver to get the correct behaviour from Inserter.Put().
var _ bigquery.ValueSaver = &Session{}

func (s *Session) Save() (map[string]bigquery.Value, string, error) {
	attributes := make([]map[string]bigquery.Value, 0, len(s.Attributes))

	for k, v := range s.Attributes {
		attributes = append(attributes, map[string]bigquery.Value{
			"key":   k,
			"value": v,
		})
	}

	row := map[string]bigquery.Value{
		"sessionId":          s.SessionID,
		"userId":             s.UserID,
		"sessionStartTime":   s.SessionStartTime.UTC(),
		"sessionEndTime":     s.SessionEndTime.UTC(),
		"applicationId":      s.ApplicationID,
		"applicationVersion": s.ApplicationVersion,
		"attributes":         attributes,
	}

	return row, s.SessionID, nil
}
