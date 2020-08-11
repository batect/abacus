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

package types_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/batect/abacus/server/types"
	"github.com/batect/abacus/server/validation"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("A session", func() {
	Describe("validation", func() {
		var v *validator.Validate
		var trans ut.Translator

		BeforeEach(func() {
			var err error
			v, trans, err = validation.CreateValidator()

			Expect(err).ToNot(HaveOccurred())
		})

		validate := func(sourceJSON string) []validation.Error {
			session := types.Session{}
			err := json.Unmarshal([]byte(sourceJSON), &session)
			Expect(err).ToNot(HaveOccurred())

			err = v.Struct(session)

			if err == nil {
				return []validation.Error{}
			}

			Expect(err).To(BeAssignableToTypeOf(validator.ValidationErrors{}))

			return validation.ToValidationErrors(err.(validator.ValidationErrors), trans)
		}

		Describe("given a valid session", func() {
			session := `{
				"sessionId": "11112222-3333-4444-5555-666677778888", 
				"userId": "99990000-3333-4444-5555-666677778888", 
				"sessionStartTime": "2019-01-02T03:04:05.678Z", 
				"sessionEndTime": "2019-01-02T09:04:05.678Z", 
				"applicationId": "test-app", 
				"applicationVersion": "1.0.0",
				"attributes": { "operatingSystem": "Mac" }
			}`

			var errors []validation.Error

			BeforeEach(func() {
				errors = validate(session)
			})

			It("returns no errors", func() {
				Expect(errors).To(BeEmpty())
			})
		})

		type invalidCase struct {
			description    string
			sourceJSON     string
			expectedErrors []validation.Error
		}

		invalidCases := []invalidCase{
			{
				description: "an empty body",
				sourceJSON:  `{}`,
				expectedErrors: []validation.Error{
					{Key: "sessionId", Type: "required", Message: "sessionId is a required field"},
					{Key: "userId", Type: "required", Message: "userId is a required field"},
					{Key: "sessionStartTime", Type: "required", Message: "sessionStartTime is a required field"},
					{Key: "sessionEndTime", Type: "required", Message: "sessionEndTime is a required field"},
					{Key: "applicationId", Type: "required", Message: "applicationId is a required field"},
					{Key: "applicationVersion", Type: "required", Message: "applicationVersion is a required field"},
				},
			},
			{
				description: "an invalid value for the ID fields",
				sourceJSON: `{
					"sessionId": "abc123", 
					"userId": "def456", 
					"sessionStartTime": "2019-01-02T03:04:05.678Z", 
					"sessionEndTime": "2019-01-02T09:04:05.678Z", 
					"applicationId": "test-app", 
					"applicationVersion": "1.0.0"
				}`,
				expectedErrors: []validation.Error{
					{Key: "sessionId", Type: "uuid", InvalidValue: "abc123", Message: "sessionId must be a valid UUID"},
					{Key: "userId", Type: "uuid", InvalidValue: "def456", Message: "userId must be a valid UUID"},
				},
			},
			{
				description: "the end time after the start time",
				sourceJSON: `{
					"sessionId": "11112222-3333-4444-5555-666677778888", 
					"userId": "99990000-3333-4444-5555-666677778888", 
					"sessionStartTime": "2019-01-04T03:04:05.678Z", 
					"sessionEndTime": "2019-01-02T09:04:05.678Z", 
					"applicationId": "test-app", 
					"applicationVersion": "1.0.0"
				}`,
				expectedErrors: []validation.Error{
					{
						Key:          "sessionEndTime",
						Type:         "gtefield",
						InvalidValue: time.Date(2019, 1, 2, 9, 4, 5, 678000000, time.UTC),
						Message:      "sessionEndTime must be greater than or equal to sessionStartTime",
					},
				},
			},
			{
				description: "an invalid application ID",
				sourceJSON: `{
					"sessionId": "11112222-3333-4444-5555-666677778888", 
					"userId": "99990000-3333-4444-5555-666677778888", 
					"sessionStartTime": "2019-01-02T03:04:05.678Z", 
					"sessionEndTime": "2019-01-02T09:04:05.678Z", 
					"applicationId": "my-app", 
					"applicationVersion": "1.0.0"
				}`,
				expectedErrors: []validation.Error{
					{Key: "applicationId", Type: "applicationId", InvalidValue: "my-app", Message: "applicationId must be a valid application ID"},
				},
			},
		}

		for _, c := range invalidCases {
			testCase := c

			Describe(fmt.Sprintf("given an invalid session with %v", testCase.description), func() {
				var errors []validation.Error

				BeforeEach(func() {
					errors = validate(testCase.sourceJSON)
				})

				It("returns the expected errors", func() {
					Expect(errors).To(ConsistOf(testCase.expectedErrors))
				})
			})
		}
	})
})
