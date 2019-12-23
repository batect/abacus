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

package api_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/batect/abacus/server/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ping endpoint", func() {
	Context("when invoked", func() {
		var resp *httptest.ResponseRecorder

		BeforeEach(func() {
			resp = httptest.NewRecorder()
			api.Ping(resp, nil)
		})

		It("returns a HTTP 200 response", func() {
			Expect(resp.Code).To(Equal(http.StatusOK))
		})

		It("returns 'pong' in the response body", func() {
			Expect(resp.Body.String()).To(Equal("pong"))
		})
	})
})
