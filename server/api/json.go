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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

const jsonMimeType = "application/json"
const contentTypeHeader = "Content-Type"

type jsonLoader struct {
	validator *validator.Validate
}

func newJSONLoader() *jsonLoader {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return &jsonLoader{
		validator: v,
	}
}

func (l *jsonLoader) LoadJSON(w http.ResponseWriter, req *http.Request, target interface{}) bool {
	if req.Header.Get(contentTypeHeader) != jsonMimeType {
		badRequest(w, "Content-Type must be 'application/json'")
		return false
	}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&target); err != nil {
		badRequest(w, fmt.Sprintf("Request body is not valid: %s", strings.TrimPrefix(err.Error(), "json: ")))
		return false
	}

	if err := l.validator.Struct(target); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			invalidBody(w, l.toValidationErrors(validationErrors))
			return false
		}

		badRequest(w, fmt.Sprintf("Request body is not valid: %s", err))

		return false
	}

	return true
}

func (l *jsonLoader) toValidationErrors(errors validator.ValidationErrors) []validationError {
	validationErrors := make([]validationError, 0, len(errors))

	for _, e := range errors {
		key := e.Namespace()

		if i := strings.Index(key, "."); i != -1 {
			key = key[i+1:]
		}

		value := fmt.Sprintf("%v", e.Value())

		if e.Tag() == "required" {
			value = ""
		}

		validationErrors = append(validationErrors, validationError{
			Key:          key,
			Type:         e.Tag(),
			InvalidValue: value,
		})
	}

	return validationErrors
}
