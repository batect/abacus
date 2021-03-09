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

package api_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/batect/abacus/server/api"
	"github.com/batect/service-observability/middleware/testutils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ping endpoint", func() {
	var resp *httptest.ResponseRecorder

	BeforeEach(func() {
		resp = httptest.NewRecorder()
	})

	Context("when invoked with a HTTP method other than GET", func() {
		BeforeEach(func() {
			req, _ := testutils.RequestWithTestLogger(httptest.NewRequest("POST", "/ping", nil))
			api.Ping(resp, req)
		})

		It("returns a HTTP 405 response", func() {
			Expect(resp.Code).To(Equal(http.StatusMethodNotAllowed))
		})

		It("returns a JSON error payload", func() {
			Expect(resp.Body).To(MatchJSON(`{"message":"This endpoint only supports GET requests"}`))
		})

		It("sets the response Content-Type header", func() {
			Expect(resp.Result().Header).To(HaveKeyWithValue("Content-Type", []string{"application/json"}))
		})

		It("sets the response Allow header", func() {
			Expect(resp.Result().Header).To(HaveKeyWithValue("Allow", []string{"GET"}))
		})
	})

	Context("when invoked with a HTTP GET", func() {
		BeforeEach(func() {
			req, _ := testutils.RequestWithTestLogger(httptest.NewRequest("GET", "/ping", nil))
			api.Ping(resp, req)
		})

		It("returns a HTTP 200 response", func() {
			Expect(resp.Code).To(Equal(http.StatusOK))
		})

		It("returns 'pong' in the response body", func() {
			Expect(resp.Body.String()).To(Equal("pong"))
		})
	})
})
