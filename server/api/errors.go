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
)

type errorResponse struct {
	Message          string            `json:"message"`
	ValidationErrors []validationError `json:"validationErrors,omitempty"`
}

type validationError struct {
	Key          string `json:"key"`
	Type         string `json:"type"`
	InvalidValue string `json:"invalidValue,omitempty"`
}

func badRequest(w http.ResponseWriter, message string) {
	resp := errorResponse{Message: message}
	resp.Write(w, http.StatusBadRequest)
}

func invalidBody(w http.ResponseWriter, errors []validationError) {
	resp := errorResponse{Message: "Request body has validation errors", ValidationErrors: errors}
	resp.Write(w, http.StatusBadRequest)
}

func methodNotAllowed(w http.ResponseWriter, allowedMethod string) {
	resp := errorResponse{Message: fmt.Sprintf("This endpoint only supports %v requests", allowedMethod)}

	w.Header().Set("Allow", allowedMethod)
	resp.Write(w, http.StatusMethodNotAllowed)
}

func (e *errorResponse) Write(w http.ResponseWriter, status int) {
	w.Header().Set(contentTypeHeader, jsonMimeType)
	w.WriteHeader(status)

	bytes, err := json.Marshal(e)

	if err != nil {
		panic(err)
	}

	if _, err := w.Write(bytes); err != nil {
		panic(err)
	}
}
