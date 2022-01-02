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

package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/batect/abacus/server/decoding"
	"github.com/batect/abacus/server/validation"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const jsonMimeType = "application/json"
const contentTypeHeader = "Content-Type"

type jsonLoader struct {
	validator  *validator.Validate
	translator ut.Translator
}

func newJSONLoader() (*jsonLoader, error) {
	v, trans, err := validation.CreateValidator()

	if err != nil {
		return nil, err
	}

	return &jsonLoader{
		validator:  v,
		translator: trans,
	}, nil
}

func (l *jsonLoader) LoadJSON(w http.ResponseWriter, req *http.Request, target interface{}) bool {
	if req.Header.Get(contentTypeHeader) != jsonMimeType {
		badRequest(req.Context(), w, "Content-Type must be 'application/json'")
		return false
	}

	decoder := decoding.NewJSONDecoder(req.Body)

	if err := decoder.Decode(&target); err != nil {
		badRequest(req.Context(), w, fmt.Sprintf("Request body is not valid: %s", strings.TrimPrefix(err.Error(), "json: ")))
		return false
	}

	if err := l.validator.Struct(target); err != nil {
		var validationErrors validator.ValidationErrors

		if errors.As(err, &validationErrors) {
			invalidBody(req.Context(), w, validation.ToValidationErrors(validationErrors, l.translator))
			return false
		}

		badRequest(req.Context(), w, fmt.Sprintf("Request body is not valid: %s", err))

		return false
	}

	return true
}
