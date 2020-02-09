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

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/option"
)

type bigQuerySessionStore struct {
	inserter *bigquery.Inserter
}

func NewBigQuerySessionStore(projectID string, datasetID string, tableID string, credsFile string) (SessionStore, error) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID, option.WithCredentialsFile(credsFile))

	if err != nil {
		return nil, err
	}

	dataset := client.Dataset(datasetID)
	table := dataset.Table(tableID)
	inserter := table.Inserter()

	return &bigQuerySessionStore{inserter}, nil
}

func (b *bigQuerySessionStore) Store(ctx context.Context, session *Session) error {
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
	metadata := make([]map[string]bigquery.Value, 0, len(s.Metadata))

	for k, v := range s.Metadata {
		metadata = append(metadata, map[string]bigquery.Value{
			"key":   k,
			"value": v,
		})
	}

	row := map[string]bigquery.Value{
		"sessionId":          s.SessionID,
		"userId":             s.UserID,
		"sessionStartTime":   s.SessionStartTime,
		"sessionEndTime":     s.SessionEndTime,
		"applicationId":      s.ApplicationID,
		"applicationVersion": s.ApplicationVersion,
		"metadata":           metadata,
	}

	return row, "", nil
}
