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

package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func loggerForRequest(logger logrus.FieldLogger, projectID string, req *http.Request) logrus.FieldLogger {
	traceID := TraceIDFromContext(req.Context())

	return logger.WithFields(logrus.Fields{
		"trace": fmt.Sprintf("projects/%s/traces/%s", projectID, traceID),
	})
}

func NewContextWithLogger(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) logrus.FieldLogger {
	return ctx.Value(loggerKey).(logrus.FieldLogger)
}

func LoggerMiddleware(baseLogger logrus.FieldLogger, projectID string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := loggerForRequest(baseLogger, projectID, req)
		logger.Debug("Processing request.")

		ctx := NewContextWithLogger(req.Context(), logger)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
