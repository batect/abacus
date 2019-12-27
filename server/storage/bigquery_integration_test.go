// Copyright 2019 Charles Korn.
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
// +build integrationTests

package storage_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/batect/abacus/server/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/api/option"
)

const testApplicationID string = "abacus-integration-tests"

var _ = AfterSuite(deleteOldTestData)

var _ = Describe("A BigQuery session store", func() {
	var store storage.SessionStore

	BeforeEach(func() {
		var err error

		store, err = storage.NewBigQuerySessionStore(getProjectID(), getDatasetID(), getSessionsTableID(), getApplicationCredentialsFilePath())

		Expect(err).ToNot(HaveOccurred())
	})

	Context("when storing a session", func() {
		var session storage.Session

		BeforeEach(func() {
			session = storage.Session{
				SessionID:          uuid.New().String(),
				UserID:             "99990000-3333-4444-5555-666677778888",
				SessionStartTime:   time.Now().Truncate(time.Microsecond).In(time.UTC),
				SessionEndTime:     time.Now().Truncate(time.Microsecond).In(time.UTC),
				ApplicationID:      testApplicationID,
				ApplicationVersion: "1.0.0",
				Metadata: map[string]string{
					"operatingSystem": "Mac",
					"dockerVersion":   "19.3.5",
				},
			}

			err := store.Store(context.Background(), &session)

			Expect(err).ToNot(HaveOccurred())
		})

		It("saves the session to BigQuery", func() {
			Expect(retrieveSession(session.SessionID)).To(Equal(session))
		})
	})
})

func deleteOldTestData() {
	client := createTestClient()

	// See https://stackoverflow.com/a/53495209/1668119 for an explanation of the condition on sessionStartTime.
	// #nosec G201
	q := client.Query(fmt.Sprintf(`
		DELETE FROM %s.%s 
		WHERE applicationID = @applicationId 
			AND sessionStartTime < TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 45 MINUTE)
		`,
		getDatasetID(),
		getSessionsTableID(),
	))

	q.Parameters = []bigquery.QueryParameter{
		{Name: "applicationId", Value: testApplicationID},
	}

	ctx := context.Background()
	job, err := q.Run(ctx)

	if err != nil {
		panic(err)
	}

	status, err := job.Wait(ctx)

	if err != nil {
		panic(err)
	}

	if status.Err() != nil {
		panic(status.Err())
	}
}

func createTestClient() *bigquery.Client {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, getProjectID(), option.WithCredentialsFile(getTestCredentialsFilePath()))

	if err != nil {
		panic(err)
	}

	return client
}

func retrieveSession(sessionID string) storage.Session {
	client := createTestClient()

	// #nosec G201
	q := client.Query(fmt.Sprintf(`
		SELECT sessionId, userId, sessionStartTime, sessionEndTime, applicationId, applicationVersion, metadata
		FROM %s.%s 
		WHERE applicationID = @applicationId
			AND sessionId = @sessionId
			AND sessionStartTime >= '2019-12-26'
		`,
		getDatasetID(),
		getSessionsTableID(),
	))

	q.Parameters = []bigquery.QueryParameter{
		{Name: "applicationId", Value: testApplicationID},
		{Name: "sessionId", Value: sessionID},
	}

	ctx := context.Background()
	it, err := q.Read(ctx)

	if err != nil {
		panic(err)
	}

	if it.TotalRows != 1 {
		panic(fmt.Sprintf("Expected to receive one row, but received %v.", it.TotalRows))
	}

	var values map[string]bigquery.Value

	if err := it.Next(&values); err != nil {
		panic(err)
	}

	return storage.Session{
		SessionID:          values["sessionId"].(string),
		UserID:             values["userId"].(string),
		SessionStartTime:   values["sessionStartTime"].(time.Time),
		SessionEndTime:     values["sessionEndTime"].(time.Time),
		ApplicationID:      values["applicationId"].(string),
		ApplicationVersion: values["applicationVersion"].(string),
		Metadata:           reconstructMetadata(values["metadata"].([]bigquery.Value)),
	}
}

func reconstructMetadata(source []bigquery.Value) map[string]string {
	metadata := map[string]string{}

	for _, m := range source {
		entry := m.(map[string]bigquery.Value)
		key := entry["key"].(string)
		value := entry["value"].(string)

		metadata[key] = value
	}

	return metadata
}

func getProjectID() string {
	return getEnvOrExit("GOOGLE_PROJECT")
}

func getDatasetID() string {
	return getEnvOrExit("DATASET_ID")
}

func getSessionsTableID() string {
	return getEnvOrExit("SESSIONS_TABLE_ID")
}

func getApplicationCredentialsFilePath() string {
	return getEnvOrExit("GOOGLE_CREDENTIALS_FILE")
}

func getTestCredentialsFilePath() string {
	return getEnvOrExit("GOOGLE_TEST_CREDENTIALS_FILE")
}

func getEnvOrExit(name string) string {
	value := os.Getenv(name)

	if value == "" {
		panic(fmt.Sprintf("Environment variable '%s' is not set.", name))
	}

	return value
}
