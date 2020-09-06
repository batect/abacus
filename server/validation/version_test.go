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

var _ = Describe("Validating versions", func() {
	var v *validator.Validate

	BeforeEach(func() {
		v = validator.New()
		en := en.New()
		uni := ut.New(en, en)
		trans, found := uni.GetTranslator("en")
		Expect(found).To(BeTrue())

		err := validation.RegisterVersionValidation(v, trans)
		Expect(err).ToNot(HaveOccurred())
	})

	type testStruct struct {
		Version string `validate:"version"`
	}

	for _, version := range []string{
		"1",
		"1.0",
		"1.0.0",
		"1.0.0-abc123",
		"1.0.0-abc-def.123",
		"1.0.0-abc123+xyz456",
		"1.0.0-abc-def.12+ghi-jkl.34",
		"1.0.0+ghi-jkl.34",
		"01.02.03",
	} {
		testObject := testStruct{version}

		Describe(fmt.Sprintf("given the valid version '%v'", testObject.Version), func() {
			It("validates as a permitted version", func() {
				Expect(v.Struct(testObject)).ToNot(HaveOccurred())
			})
		})
	}

	Describe("given an empty version", func() {
		testObject := testStruct{Version: ""}

		It("fails validation", func() {
			Expect(v.Struct(testObject)).To(MatchError("Key: 'testStruct.Version' Error:Field validation for 'Version' failed on the 'version' tag"))
		})
	})

	for _, version := range []string{
		"1.",
		"1.0.",
		"1.0.0.",
		"1.2.3.4",
		"1.0.0-abc123+",
		"-1.2.3",
		"a",
		"a.b.c",
		".0.1",
		"1-",
		"1-thing",
		"1.2-",
		"1.2-thing",
		"1.2.3-",
		"1..2",
		"1.0.0-/../",
		"1.0.0-/..",
		"1.0.0-/",
		"1.0.0-../",
		"1.0.0+/../",
		"1.0.0+/..",
		"1.0.0+/",
		"1.0.0+../",
	} {
		testObject := testStruct{version}

		Describe(fmt.Sprintf("given the invalid version '%v'", testObject.Version), func() {
			It("fails validation", func() {
				Expect(v.Struct(testObject)).To(MatchError("Key: 'testStruct.Version' Error:Field validation for 'Version' failed on the 'version' tag"))
			})
		})
	}
})
