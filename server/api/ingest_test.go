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
	currentTime := time.Date(2020, 5, 24, 10, 12, 14, 123, time.UTC)

	BeforeEach(func() {
		store = &mockStore{}
		timeSource := func() time.Time { return currentTime }

		var err error
		handler, err = api.NewIngestHandlerWithTimeSource(store, timeSource)
		Expect(err).ToNot(HaveOccurred())

		resp = httptest.NewRecorder()
	})

	ItReturnsABadRequestResponseWithBody := func(expectedBody string) {
		It("returns a HTTP 400 response", func() {
			Expect(resp.Code).To(Equal(http.StatusBadRequest))
		})

		It("returns a JSON error payload", func() {
			Expect(resp.Body).To(MatchJSON(expectedBody))
		})

		It("sets the response Content-Type header", func() {
			Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
		})

		It("does not store any sessions", func() {
			Expect(store.StoredSessions).To(BeEmpty())
		})
	}

	Context("when invoked with a HTTP method other than PUT", func() {
		BeforeEach(func() {
			req, _ := testutils.RequestWithTestLogger(httptest.NewRequest("POST", "/ingest", nil))
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
				req, _ := testutils.RequestWithTestLogger(httptest.NewRequest("PUT", "/ingest", nil))
				handler.ServeHTTP(resp, req)
			})

			ItReturnsABadRequestResponseWithBody(`{"message":"Content-Type must be 'application/json'"}`)
		})

		Context("when invoked with an invalid Content-Type header", func() {
			BeforeEach(func() {
				req, _ := testutils.RequestWithTestLogger(httptest.NewRequest("PUT", "/ingest", nil))
				req.Header.Set("Content-Type", "text/plain")
				handler.ServeHTTP(resp, req)
			})

			ItReturnsABadRequestResponseWithBody(`{"message":"Content-Type must be 'application/json'"}`)
		})

		Context("when invoked with the required Content-Type header", func() {
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

				ItReturnsABadRequestResponseWithBody(`{"message":"Request body is not valid: EOF"}`)
			})

			Context("when the request body is not valid JSON", func() {
				BeforeEach(func() {
					req, _ := createRequest("{")
					handler.ServeHTTP(resp, req)
				})

				ItReturnsABadRequestResponseWithBody(`{"message":"Request body is not valid: unexpected EOF"}`)
			})

			Context("when the request body is valid JSON but is empty", func() {
				BeforeEach(func() {
					req, _ := createRequest("{}")
					handler.ServeHTTP(resp, req)
				})

				ItReturnsABadRequestResponseWithBody(`{
					"message": "Request body has validation errors",
					"validationErrors": [
						{ "key": "sessionId", "type": "required", "message": "sessionId is a required field" },
						{ "key": "userId", "type": "required", "message": "userId is a required field" },
						{ "key": "sessionStartTime", "type": "required", "message": "sessionStartTime is a required field" },
						{ "key": "sessionEndTime", "type": "required", "message": "sessionEndTime is a required field" },
						{ "key": "applicationId", "type": "required", "message": "applicationId is a required field" },
						{ "key": "applicationVersion", "type": "required", "message": "applicationVersion is a required field" }
					]
				}`)
			})

			Context("when the request body is valid JSON but has an invalid value for one or more fields", func() {
				Context("because one or more fields is in an invalid format", func() {
					BeforeEach(func() {
						req, _ := createRequest(`{
							"sessionId": "abc123", 
							"userId": "def456", 
							"sessionStartTime": "2019-01-02T03:04:05.678Z", 
							"sessionEndTime": "2019-01-02T09:04:05.678Z", 
							"applicationId": "test-app", 
							"applicationVersion": "1.0.0"
						}`)

						handler.ServeHTTP(resp, req)
					})

					ItReturnsABadRequestResponseWithBody(`{
						"message": "Request body has validation errors",
						"validationErrors": [
							{ "key": "sessionId", "type": "uuid", "invalidValue": "abc123", "message": "sessionId must be a valid UUID" },
							{ "key": "userId", "type": "uuid", "invalidValue": "def456", "message": "userId must be a valid UUID" }
						]
					}`)
				})

				Context("because the session start time is after the session end time", func() {
					BeforeEach(func() {
						req, _ := createRequest(`{
							"sessionId": "11112222-3333-4444-5555-666677778888", 
							"userId": "99990000-3333-4444-5555-666677778888", 
							"sessionStartTime": "2019-01-04T03:04:05.678Z", 
							"sessionEndTime": "2019-01-02T09:04:05.678Z", 
							"applicationId": "test-app", 
							"applicationVersion": "1.0.0"
						}`)

						handler.ServeHTTP(resp, req)
					})

					ItReturnsABadRequestResponseWithBody(`{
						"message": "Request body has validation errors",
						"validationErrors": [
							{
								"key": "sessionEndTime", 
								"type": "gtefield", 
								"invalidValue": "2019-01-02T09:04:05.678Z", 
								"message": "sessionEndTime must be greater than or equal to sessionStartTime"
							}
						]
					}`)
				})

				Context("because the application ID is invalid", func() {
					BeforeEach(func() {
						req, _ := createRequest(`{
							"sessionId": "11112222-3333-4444-5555-666677778888", 
							"userId": "99990000-3333-4444-5555-666677778888", 
							"sessionStartTime": "2019-01-02T03:04:05.678Z", 
							"sessionEndTime": "2019-01-02T09:04:05.678Z", 
							"applicationId": "my-app", 
							"applicationVersion": "1.0.0"
						}`)

						handler.ServeHTTP(resp, req)
					})

					ItReturnsABadRequestResponseWithBody(`{
						"message": "Request body has validation errors",
						"validationErrors": [
							{ "key": "applicationId", "type": "applicationId", "invalidValue": "my-app", "message": "applicationId must be a valid application ID" }
						]
					}`)
				})
			})

			Context("when the request body is valid JSON but has an extra field", func() {
				BeforeEach(func() {
					req, _ := createRequest(`{"sessionId": "11112222-3333-4444-5555-666677778888", "blah": "value"}`)
					handler.ServeHTTP(resp, req)
				})

				ItReturnsABadRequestResponseWithBody(`{"message":"Request body is not valid: unknown field \"blah\""}`)
			})

			Context("when the request body is valid JSON but contains a value for the ingestion time", func() {
				BeforeEach(func() {
					req, _ := createRequest(`{"sessionId": "11112222-3333-4444-5555-666677778888", "IngestionTime": "2020-05-24T00:00:00.000Z"}`)
					handler.ServeHTTP(resp, req)
				})

				ItReturnsABadRequestResponseWithBody(`{"message":"Request body is not valid: unknown field \"IngestionTime\""}`)
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
						"applicationId": "test-app", 
						"applicationVersion": "1.0.0"
					}`)
				})

				Context("when the session does not already exist", func() {
					BeforeEach(func() {
						store.SessionExists = false
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
								IngestionTime:      currentTime,
								ApplicationID:      "test-app",
								ApplicationVersion: "1.0.0",
							}))
						})
					})

					Context("when storing the session fails", func() {
						BeforeEach(func() {
							store.ErrorToReturnFromStore = errors.New("could not store session")
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
							Expect(loggingHook.Entries).To(ContainElement(LogEntryWithError("Storing session failed.", store.ErrorToReturnFromStore)))
						})
					})
				})

				Context("when the session already exists", func() {
					BeforeEach(func() {
						store.SessionExists = true
						handler.ServeHTTP(resp, req)
					})

					It("returns a HTTP 304 response", func() {
						Expect(resp.Code).To(Equal(http.StatusNotModified))
					})

					It("returns an empty body", func() {
						Expect(resp.Result().ContentLength).To(BeZero())
					})

					It("does not store the session", func() {
						Expect(store.StoredSessions).To(BeEmpty())
					})

					It("logs a warning", func() {
						Expect(loggingHook.Entries).To(ConsistOf(LogEntryWithWarning("Session already exists, not storing.")))
					})
				})

				Context("when checking if the session already exists fails", func() {
					BeforeEach(func() {
						store.ErrorToReturnFromCheckIfExists = errors.New("could not check if session exists")
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
						Expect(loggingHook.Entries).To(ContainElement(LogEntryWithError("Checking if session already exists failed.", store.ErrorToReturnFromCheckIfExists)))
					})

					It("does not store the session", func() {
						Expect(store.StoredSessions).To(BeEmpty())
					})
				})
			})
		})
	})
})

type mockStore struct {
	ErrorToReturnFromStore         error
	StoredSessions                 []storage.Session
	ErrorToReturnFromCheckIfExists error
	SessionExists                  bool
}

func (m *mockStore) Store(_ context.Context, session *storage.Session) error {
	if m.ErrorToReturnFromStore != nil {
		return m.ErrorToReturnFromStore
	}

	m.StoredSessions = append(m.StoredSessions, *session)

	return nil
}

func (m *mockStore) CheckIfExists(_ context.Context, session *storage.Session) (bool, error) {
	return m.SessionExists, m.ErrorToReturnFromCheckIfExists
}

func GetMessage(e logrus.Entry) string     { return e.Message }
func GetData(e logrus.Entry) logrus.Fields { return e.Data }
func GetLevel(e logrus.Entry) logrus.Level { return e.Level }

func LogEntryWithError(message string, err error) types.GomegaMatcher {
	return SatisfyAll(
		WithTransform(GetMessage, Equal(message)),
		WithTransform(GetData, Equal(logrus.Fields{logrus.ErrorKey: err})),
		WithTransform(GetLevel, Equal(logrus.ErrorLevel)),
	)
}

func LogEntryWithWarning(message string) types.GomegaMatcher {
	return SatisfyAll(
		WithTransform(GetMessage, Equal(message)),
		WithTransform(GetLevel, Equal(logrus.WarnLevel)),
	)
}
