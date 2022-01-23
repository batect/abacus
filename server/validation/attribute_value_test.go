// Copyright 2019-2022 Charles Korn.
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
	"encoding/json"
	"fmt"

	"github.com/batect/abacus/server/validation"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validating attribute values", func() {
	var v *validator.Validate

	BeforeEach(func() {
		v = validator.New()
		en := en.New()
		uni := ut.New(en, en)
		trans, found := uni.GetTranslator("en")
		Expect(found).To(BeTrue())

		err := validation.RegisterAttributeValueValidation(v, trans)
		Expect(err).ToNot(HaveOccurred())
	})

	type testStruct struct {
		AttributeValue interface{} `validate:"attributeValue"`
	}

	for _, id := range []interface{}{
		"a",
		"abc123",
		"123",
		json.Number("0"),
		json.Number("-1"),
		json.Number("1"),
		json.Number("1.5"),
		true,
		false,
		nil,
	} {
		testObject := testStruct{id}

		Describe(fmt.Sprintf("given the attribute value %#v", testObject.AttributeValue), func() {
			It("validates as a permitted attribute value", func() {
				Expect(v.Struct(testObject)).ToNot(HaveOccurred())
			})
		})
	}

	for _, value := range []interface{}{
		[]string{},
		[]string{"hello"},
		map[string]string{
			"hello": "world",
		},
	} {
		testObject := testStruct{value}

		Describe(fmt.Sprintf("given the attribute value %#v", testObject.AttributeValue), func() {
			It("fails validation", func() {
				Expect(v.Struct(testObject)).To(MatchError("Key: 'testStruct.AttributeValue' Error:Field validation for 'AttributeValue' failed on the 'attributeValue' tag"))
			})
		})
	}
})
