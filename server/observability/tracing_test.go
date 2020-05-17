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

package observability_test

import (
	"net/http/httptest"

	"github.com/batect/abacus/server/observability"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Naming HTTP request spans", func() {
	req := httptest.NewRequest("PUT", "/blah", nil)

	Describe("given an operation name is provided", func() {
		var name string

		BeforeEach(func() {
			name = observability.NameHTTPRequestSpan("Server", req)
		})

		It("includes the operation name, HTTP method and URL in the name", func() {
			Expect(name).To(Equal("Server: PUT /blah"))
		})
	})

	Describe("given an operation name is not provided", func() {
		var name string

		BeforeEach(func() {
			name = observability.NameHTTPRequestSpan("", req)
		})

		It("does not include the operation name in the span name", func() {
			Expect(name).To(Equal("PUT /blah"))
		})
	})
})
