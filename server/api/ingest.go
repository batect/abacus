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
	"time"

	"github.com/batect/abacus/server/storage"
	"github.com/batect/abacus/server/types"
	"github.com/batect/service-observability/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ingestHandler struct {
	loader       *jsonLoader
	sessionStore storage.SessionStore
	timeSource   timeSource
}

type timeSource func() time.Time

const sessionID = attribute.Key("session.sessionId")
const userID = attribute.Key("session.userId")
const applicationID = attribute.Key("session.applicationId")
const applicationVersion = attribute.Key("session.applicationVersion")

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

	session := types.Session{}

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
		userID.String(session.UserID),
		applicationID.String(session.ApplicationID),
		applicationVersion.String(session.ApplicationVersion),
	)

	session = h.cleanSession(session)

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

func (h *ingestHandler) cleanSession(session types.Session) types.Session {
	session.IngestionTime = h.timeSource()

	if session.Attributes == nil {
		session.Attributes = map[string]interface{}{}
	}

	if session.Events == nil {
		session.Events = []types.Event{}
	}

	if session.Spans == nil {
		session.Spans = []types.Span{}
	}

	for i, s := range session.Spans {
		if s.Attributes == nil {
			session.Spans[i].Attributes = map[string]interface{}{}
		}
	}

	for i, e := range session.Events {
		if e.Attributes == nil {
			session.Events[i].Attributes = map[string]interface{}{}
		}
	}

	return session
}
