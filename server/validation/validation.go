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

package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Error struct {
	Key          string      `json:"key"`
	Type         string      `json:"type"`
	InvalidValue interface{} `json:"invalidValue,omitempty"`
	Message      string      `json:"message"`
}

func CreateValidator() (*validator.Validate, ut.Translator, error) {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	en := en.New()
	uni := ut.New(en, en)
	trans, found := uni.GetTranslator("en")

	if !found {
		return nil, nil, errors.New("could not load English translator")
	}

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, nil, fmt.Errorf("could not register default translations: %w", err)
	}

	if err := RegisterApplicationIDValidation(v, trans); err != nil {
		return nil, nil, fmt.Errorf("could not register application ID validator: %w", err)
	}

	if err := RegisterVersionValidation(v, trans); err != nil {
		return nil, nil, fmt.Errorf("could not register version validator: %w", err)
	}

	return v, trans, nil
}

func ToValidationErrors(errors validator.ValidationErrors, translator ut.Translator) []Error {
	validationErrors := make([]Error, 0, len(errors))

	for _, e := range errors {
		key := e.Namespace()

		if i := strings.Index(key, "."); i != -1 {
			key = key[i+1:]
		}

		value := e.Value()

		if e.Tag() == "required" {
			value = nil
		}

		validationErrors = append(validationErrors, Error{
			Key:          key,
			Type:         e.Tag(),
			InvalidValue: value,
			Message:      e.Translate(translator),
		})
	}

	return validationErrors
}


func registerValidation(v *validator.Validate, trans ut.Translator, tag string, errorMessage string, validationFunc validator.Func) error {
	if err := v.RegisterValidation(tag, validationFunc, false); err != nil {
		return fmt.Errorf("could not register %v validator: %w", tag, err)
	}

	if err := v.RegisterTranslation(tag, trans, registrationFunc(tag, errorMessage), translateFunc); err != nil {
		return fmt.Errorf("could not register %v validator error message translation: %w", tag, err)
	}

	return nil
}

func registrationFunc(tag string, translation string) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		if err := ut.Add(tag, translation, false); err != nil {
			return err
		}

		return nil
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())

	if err != nil {
		panic(fmt.Sprintf("error translating FieldError: %#v", fe))
	}

	return t
}
