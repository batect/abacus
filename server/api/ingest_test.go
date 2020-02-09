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

package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/batect/abacus/server/api"
	"github.com/batect/abacus/server/middleware/testutils"
	"github.com/batect/abacus/server/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

var _ = Describe("Ingest endpoint", func() {
	var handler http.Handler
	var resp *httptest.ResponseRecorder
	var store *mockStore

	BeforeEach(func() {
		store = &mockStore{}
		handler = api.NewIngestHandler(store)
		resp = httptest.NewRecorder()
	})

	Context("when invoked with a HTTP method other than PUT", func() {
		BeforeEach(func() {
			req := httptest.NewRequest("POST", "/ingest", nil)
			handler.ServeHTTP(resp, req)
		})

		It("returns a HTTP 405 response", func() {
			Expect(resp.Code).To(Equal(http.StatusMethodNotAllowed))
		})

		It("returns a JSON error payload", func() {
			Expect(resp.Body).To(MatchJSON(`{"message":"This endpoint only supports PUT requests"}`))
		})

		It("sets the response Content-Type header", func() {
			Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
		})

		It("sets the response Allow header", func() {
			Expect(resp.Result().Header).To(HaveKeyWithValue("Allow", []string{"PUT"}))
		})

		It("does not store any sessions", func() {
			Expect(store.StoredSessions).To(BeEmpty())
		})
	})

	Context("when invoked with PUT", func() {
		Context("when invoked with no Content-Type header", func() {
			BeforeEach(func() {
				req := httptest.NewRequest("PUT", "/ingest", nil)
				handler.ServeHTTP(resp, req)
			})

			It("returns a HTTP 400 response", func() {
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})

			It("returns a JSON error payload", func() {
				Expect(resp.Body).To(MatchJSON(`{"message":"Content-Type must be 'application/json'"}`))
			})

			It("sets the response Content-Type header", func() {
				Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
			})

			It("does not store any sessions", func() {
				Expect(store.StoredSessions).To(BeEmpty())
			})
		})

		Context("when invoked with an invalid Content-Type header", func() {
			BeforeEach(func() {
				req := httptest.NewRequest("PUT", "/ingest", nil)
				req.Header.Set("Content-Type", "text/plain")
				handler.ServeHTTP(resp, req)
			})

			It("returns a HTTP 400 response", func() {
				Expect(resp.Code).To(Equal(http.StatusBadRequest))
			})

			It("returns a JSON error payload", func() {
				Expect(resp.Body).To(MatchJSON(`{"message":"Content-Type must be 'application/json'"}`))
			})

			It("sets the response Content-Type header", func() {
				Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
			})

			It("does not store any sessions", func() {
				Expect(store.StoredSessions).To(BeEmpty())
			})
		})

		Context("when invoked with a the required Content-Type header", func() {
			createRequest := func(body string) (*http.Request, *test.Hook) {
				req := httptest.NewRequest("PUT", "/ingest", strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				return testutils.RequestWithTestLogger(req)
			}

			Context("when the request body is empty", func() {
				BeforeEach(func() {
					req, _ := createRequest("")
					handler.ServeHTTP(resp, req)
				})

				It("returns a HTTP 400 response", func() {
					Expect(resp.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns a JSON error payload", func() {
					Expect(resp.Body).To(MatchJSON(`{"message": "Request body is not valid: EOF"}`))
				})

				It("sets the response Content-Type header", func() {
					Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
				})

				It("does not store any sessions", func() {
					Expect(store.StoredSessions).To(BeEmpty())
				})
			})

			Context("when the request body is not valid JSON", func() {
				BeforeEach(func() {
					req, _ := createRequest("{")
					handler.ServeHTTP(resp, req)
				})

				It("returns a HTTP 400 response", func() {
					Expect(resp.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns a JSON error payload", func() {
					Expect(resp.Body).To(MatchJSON(`{"message": "Request body is not valid: unexpected EOF"}`))
				})

				It("sets the response Content-Type header", func() {
					Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
				})

				It("does not store any sessions", func() {
					Expect(store.StoredSessions).To(BeEmpty())
				})
			})

			Context("when the request body is valid JSON but is empty", func() {
				BeforeEach(func() {
					req, _ := createRequest("{}")
					handler.ServeHTTP(resp, req)
				})

				It("returns a HTTP 400 response", func() {
					Expect(resp.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns a JSON error payload", func() {
					Expect(resp.Body).To(MatchJSON(`{
						"message": "Request body has validation errors",
						"validationErrors": [
							{ "key": "sessionId", "type": "required" },
							{ "key": "userId", "type": "required" },
							{ "key": "sessionStartTime", "type": "required" },
							{ "key": "sessionEndTime", "type": "required" },
							{ "key": "applicationId", "type": "required" },
							{ "key": "applicationVersion", "type": "required" }
						]
					}`))
				})

				It("sets the response Content-Type header", func() {
					Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
				})

				It("does not store any sessions", func() {
					Expect(store.StoredSessions).To(BeEmpty())
				})
			})

			Context("when the request body is valid JSON but has an invalid value for one or more fields", func() {
				BeforeEach(func() {
					req, _ := createRequest(`{
						"sessionId": "abc123", 
						"userId": "def456", 
						"sessionStartTime": "2019-01-02T03:04:05.678Z", 
						"sessionEndTime": "2019-01-02T09:04:05.678Z", 
						"applicationId": "my-app", 
						"applicationVersion": "1.0.0"
					}`)

					handler.ServeHTTP(resp, req)
				})

				It("returns a HTTP 400 response", func() {
					Expect(resp.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns a JSON error payload with details of each of the errors", func() {
					Expect(resp.Body).To(MatchJSON(`{
						"message": "Request body has validation errors",
						"validationErrors": [
							{ "key": "sessionId", "type": "uuid", "invalidValue": "abc123" },
							{ "key": "userId", "type": "uuid", "invalidValue": "def456" }
						]
					}`))
				})

				It("sets the response Content-Type header", func() {
					Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
				})

				It("does not store any sessions", func() {
					Expect(store.StoredSessions).To(BeEmpty())
				})
			})

			Context("when the request body is valid JSON but has an extra field", func() {
				BeforeEach(func() {
					req, _ := createRequest(`{"sessionId": "11112222-3333-4444-5555-666677778888", "blah": "value"}`)
					handler.ServeHTTP(resp, req)
				})

				It("returns a HTTP 400 response", func() {
					Expect(resp.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns a JSON error payload", func() {
					Expect(resp.Body).To(MatchJSON(`{"message": "Request body is not valid: unknown field \"blah\""}`))
				})

				It("sets the response Content-Type header", func() {
					Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
				})

				It("does not store any sessions", func() {
					Expect(store.StoredSessions).To(BeEmpty())
				})
			})

			Context("when the request body is valid", func() {
				var req *http.Request
				var loggingHook *test.Hook

				BeforeEach(func() {
					req, loggingHook = createRequest(`{
						"sessionId": "11112222-3333-4444-5555-666677778888", 
						"userId": "99990000-3333-4444-5555-666677778888", 
						"sessionStartTime": "2019-01-02T03:04:05.678Z", 
						"sessionEndTime": "2019-01-02T09:04:05.678Z", 
						"applicationId": 
						"my-app", 
						"applicationVersion": "1.0.0"
					}`)
				})

				Context("when storing the session succeeds", func() {
					BeforeEach(func() {
						handler.ServeHTTP(resp, req)
					})

					It("returns a HTTP 201 response", func() {
						Expect(resp.Code).To(Equal(http.StatusCreated))
					})

					It("returns an empty body", func() {
						Expect(resp.Result().ContentLength).To(BeZero())
					})

					It("stores the session", func() {
						Expect(store.StoredSessions).To(ConsistOf(storage.Session{
							SessionID:          "11112222-3333-4444-5555-666677778888",
							UserID:             "99990000-3333-4444-5555-666677778888",
							SessionStartTime:   time.Date(2019, 1, 2, 3, 4, 5, 678000000, time.UTC),
							SessionEndTime:     time.Date(2019, 1, 2, 9, 4, 5, 678000000, time.UTC),
							ApplicationID:      "my-app",
							ApplicationVersion: "1.0.0",
						}))
					})
				})

				Context("when storing the session fails", func() {
					BeforeEach(func() {
						store.ErrorToReturn = errors.New("could not store session")
						handler.ServeHTTP(resp, req)
					})

					It("returns a HTTP 503 response", func() {
						Expect(resp.Code).To(Equal(http.StatusServiceUnavailable))
					})

					It("returns a JSON error payload", func() {
						Expect(resp.Body).To(MatchJSON(`{"message": "Could not process request"}`))
					})

					It("sets the response Content-Type header", func() {
						Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
					})

					It("logs the error", func() {
						Expect(loggingHook.Entries).To(ConsistOf(LogEntryWithError("Storing session failed", store.ErrorToReturn)))
					})
				})
			})
		})
	})
})

type mockStore struct {
	ErrorToReturn  error
	StoredSessions []storage.Session
}

func (m *mockStore) Store(_ context.Context, session *storage.Session) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	m.StoredSessions = append(m.StoredSessions, *session)

	return nil
}

func LogEntryWithError(message string, err error) types.GomegaMatcher {
	GetMessage := func(e logrus.Entry) string { return e.Message }
	GetData := func(e logrus.Entry) logrus.Fields { return e.Data }
	GetLevel := func(e logrus.Entry) logrus.Level { return e.Level }

	return SatisfyAll(
		WithTransform(GetMessage, Equal(message)),
		WithTransform(GetData, Equal(logrus.Fields{logrus.ErrorKey: err})),
		WithTransform(GetLevel, Equal(logrus.ErrorLevel)),
	)
}
