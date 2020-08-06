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
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/batect/abacus/server/middleware"
	"github.com/batect/abacus/server/storage"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
)

type ingestHandler struct {
	loader       *jsonLoader
	sessionStore storage.SessionStore
	timeSource   timeSource
}

type timeSource func() time.Time

const sessionID kv.Key = kv.Key("sessionId")
const applicationID kv.Key = kv.Key("applicationId")

func NewIngestHandler(sessionStore storage.SessionStore) (http.Handler, error) {
	return NewIngestHandlerWithTimeSource(sessionStore, time.Now)
}

func NewIngestHandlerWithTimeSource(sessionStore storage.SessionStore, timeSource timeSource) (http.Handler, error) {
	loader, err := newJSONLoader()

	if err != nil {
		return nil, fmt.Errorf("could not create JSON loader: %w", err)
	}

	return &ingestHandler{
		loader:       loader,
		sessionStore: sessionStore,
		timeSource:   timeSource,
	}, nil
}

func (h *ingestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !requireMethod(w, req, http.MethodPut) {
		return
	}

	session := storage.Session{}

	if ok := h.loader.LoadJSON(w, req, &session); !ok {
		return
	}

	log := middleware.LoggerFromContext(req.Context()).
		WithField("sessionId", session.SessionID).
		WithField("applicationId", session.ApplicationID)

	ctx := middleware.ContextWithLogger(req.Context(), log)

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		sessionID.String(session.SessionID),
		applicationID.String(session.ApplicationID),
	)

	session.IngestionTime = h.timeSource()

	if err := h.sessionStore.Store(ctx, &session); errors.Is(err, storage.ErrAlreadyExists) {
		log.Warn("Session already exists, not storing.")

		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusNotModified)

		return
	} else if err != nil {
		log.WithError(err).Error("Storing session failed.")

		resp := errorResponse{Message: "Could not process request"}
		resp.Write(ctx, w, http.StatusServiceUnavailable)

		return
	}

	log.Info("Stored session successfully.")

	w.Header().Set("Content-Length", "0")
	w.WriteHeader(http.StatusCreated)
}
