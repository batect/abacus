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
// +build unitTests

package storage

import (
	"time"

	"cloud.google.com/go/bigquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("A session", func() {
	session := Session{
		SessionID:          "11112222-3333-4444-5555-666677778888",
		UserID:             "99990000-3333-4444-5555-666677778888",
		SessionStartTime:   time.Date(2019, 1, 2, 3, 4, 5, 678000000, time.UTC),
		SessionEndTime:     time.Date(2019, 1, 2, 9, 4, 5, 678000000, time.UTC),
		ApplicationID:      "my-app",
		ApplicationVersion: "1.0.0",
		Metadata: map[string]string{
			"operatingSystem": "Mac",
			"dockerVersion":   "19.3.5",
		},
	}

	Context("when stored in BigQuery", func() {
		row, insertID, err := session.Save()

		It("does not return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("does not return an insert ID", func() {
			Expect(insertID).To(BeEmpty())
		})

		It("converts the session to the format expected by BigQuery", func() {
			Expect(row).To(Equal(map[string]bigquery.Value{
				"sessionId":          session.SessionID,
				"userId":             session.UserID,
				"sessionStartTime":   "2019-01-02T03:04:05.678Z",
				"sessionEndTime":     "2019-01-02T09:04:05.678Z",
				"applicationId":      session.ApplicationID,
				"applicationVersion": session.ApplicationVersion,
				"metadata": []map[string]bigquery.Value{
					{
						"key":   "operatingSystem",
						"value": "Mac",
					},
					{
						"key":   "dockerVersion",
						"value": "19.3.5",
					},
				},
			}))
		})
	})
})
