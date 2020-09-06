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

package validation_test

import (
	"fmt"

	"github.com/batect/abacus/server/validation"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validating attribute names", func() {
	var v *validator.Validate

	BeforeEach(func() {
		v = validator.New()
		en := en.New()
		uni := ut.New(en, en)
		trans, found := uni.GetTranslator("en")
		Expect(found).To(BeTrue())

		err := validation.RegisterAttributeNameValidation(v, trans)
		Expect(err).ToNot(HaveOccurred())
	})

	type testStruct struct {
		AttributeName string `validate:"attributeName"`
	}

	for _, id := range []string{"a", "abc123", "theThing"} {
		testObject := testStruct{id}

		Describe(fmt.Sprintf("given the attribute name '%v'", testObject.AttributeName), func() {
			It("validates as a permitted attribute name", func() {
				Expect(v.Struct(testObject)).ToNot(HaveOccurred())
			})
		})
	}

	Describe("given an empty attribute name", func() {
		testObject := testStruct{AttributeName: ""}

		It("fails validation", func() {
			Expect(v.Struct(testObject)).To(MatchError("Key: 'testStruct.AttributeName' Error:Field validation for 'AttributeName' failed on the 'attributeName' tag"))
		})
	})

	for _, name := range []string{"", "1", "1a", "-", ".", "("} {
		testObject := testStruct{name}

		Describe(fmt.Sprintf("given the attribute name '%v'", testObject.AttributeName), func() {
			It("fails validation", func() {
				Expect(v.Struct(testObject)).To(MatchError("Key: 'testStruct.AttributeName' Error:Field validation for 'AttributeName' failed on the 'attributeName' tag"))
			})
		})
	}
})
