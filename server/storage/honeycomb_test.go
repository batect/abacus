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
// +build unitTests

package storage_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/batect/abacus/server/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Honeycomb session storage", func() {
	session := storage.Session{
		SessionID:          "11112222-3333-4444-5555-666677778888",
		UserID:             "99990000-3333-4444-5555-666677778888",
		SessionStartTime:   time.Date(2019, 1, 2, 3, 4, 5, 678000001, time.UTC),
		SessionEndTime:     time.Date(2019, 1, 2, 9, 4, 5, 678000000, time.UTC),
		ApplicationID:      "my-app",
		ApplicationVersion: "1.0.0",
		Metadata: map[string]string{
			"operatingSystem": "Mac",
			"dockerVersion":   "19.3.5",
		},
	}

	var testServer *httptest.Server
	var testServerURL url.URL
	var store storage.SessionStore
	var requestSent http.Request
	var requestBodySent string
	var responseCode int
	var responseBody string
	var responseDelay time.Duration

	BeforeEach(func() {
		responseCode = http.StatusTeapot
		responseBody = "Not set in test!"
		responseDelay = 0

		testServer = httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			requestSent = *request
			requestBodySent = readAllBytes(request.Body)

			time.Sleep(responseDelay)

			writer.WriteHeader(responseCode)

			if _, err := writer.Write([]byte(responseBody)); err != nil {
				panic(err)
			}
		}))

		testServerURL = mustParseURL(testServer.URL)

		store = storage.NewHoneycombSessionStore(testServerURL, "my-dataset", "my-api-key")
	})

	AfterEach(func() {
		testServer.Close()
	})

	Context("when sending the session to Honeycomb succeeds", func() {
		var err error

		BeforeEach(func() {
			responseCode = http.StatusOK

			err = store.Store(context.Background(), &session)
		})

		It("does not return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("sends the request to the configured Honeycomb hostname", func() {
			Expect(requestSent.Host).To(Equal(testServerURL.Host))
		})

		It("uses the correct URL, including the dataset name", func() {
			Expect(requestSent.URL.String()).To(Equal("/1/events/my-dataset"))
		})

		It("sends a POST request", func() {
			Expect(requestSent.Method).To(Equal(http.MethodPost))
		})

		It("includes the API key in the request headers", func() {
			Expect(requestSent.Header).To(HaveKeyWithValue("X-Honeycomb-Team", []string{"my-api-key"}))
		})

		It("includes the session start time as the Honeycomb event time", func() {
			Expect(requestSent.Header).To(HaveKeyWithValue("X-Honeycomb-Event-Time", []string{"2019-01-02T03:04:05.678000001Z"}))
		})

		It("includes a Content-Type header in the request", func() {
			Expect(requestSent.Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
		})

		It("includes the session in the request body", func() {
			Expect(requestBodySent).To(MatchJSON(`
				{
					"sessionId": "11112222-3333-4444-5555-666677778888",
					"userId": "99990000-3333-4444-5555-666677778888",
					"sessionStartTime": "2019-01-02T03:04:05.678000001Z",
					"sessionEndTime": "2019-01-02T09:04:05.678Z",
					"applicationId": "my-app",
					"applicationVersion": "1.0.0",
					"metadata": {
						"operatingSystem": "Mac",
						"dockerVersion": "19.3.5"
					}
				}
			`))
		})
	})

	Context("when sending the session to Honeycomb fails", func() {
		var err error

		BeforeEach(func() {
			responseCode = http.StatusBadRequest
			responseBody = "malformed request"

			err = store.Store(context.Background(), &session)
		})

		It("returns an error", func() {
			Expect(err).To(MatchError("request to Honeycomb (POST " + testServerURL.String() + "/1/events/my-dataset) received error response: HTTP 400: malformed request"))
		})
	})

	Context("when Honeycomb does not respond within the expected time", func() {
		var err error

		BeforeEach(func() {
			responseCode = http.StatusOK
			responseDelay = 3 * time.Second

			err = store.Store(context.Background(), &session)
		})

		It("returns an error", func() {
			Expect(err).To(MatchError(
				"request to Honeycomb (POST " + testServerURL.String() +
					"/1/events/my-dataset) failed with error: Post \"" + testServerURL.String() +
					"/1/events/my-dataset\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)",
			))
		})
	})
})

func mustParseURL(input string) url.URL {
	var parsed *url.URL
	var err error

	if parsed, err = url.Parse(input); err != nil {
		panic(err)
	}

	return *parsed
}

func readAllBytes(reader io.Reader) string {
	bytes, err := ioutil.ReadAll(reader)

	if err != nil {
		panic(err)
	}

	return string(bytes)
}
