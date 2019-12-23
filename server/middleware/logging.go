// Copyright 2019 Charles Korn.
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
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type contextKey int

const (
	loggerKey contextKey = iota
)

// Based on https://github.com/TV4/logrus-stackdriver-formatter#http-request-context and
// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#HttpRequest
func loggerForRequest(logger logrus.FieldLogger, req *http.Request) logrus.FieldLogger {
	remoteIP := req.RemoteAddr

	if strings.Contains(remoteIP, ":") {
		remoteIP = remoteIP[:strings.LastIndex(remoteIP, ":")]
	}

	return logger.WithFields(logrus.Fields{
		"httpRequest": map[string]interface{}{
			"requestMethod": req.Method,
			"requestUrl":    req.URL.String(),
			"requestSize":   req.ContentLength,
			"userAgent":     req.Header.Get("User-Agent"),
			"remoteIp":      remoteIP,
			"referrer":      req.Header.Get("Referer"),
			"protocol":      req.Proto,
		},
	})
}

func newContextWithLogger(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) logrus.FieldLogger {
	return ctx.Value(loggerKey).(logrus.FieldLogger)
}

func LoggerMiddleware(baseLogger logrus.FieldLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := loggerForRequest(baseLogger, req)
		logger.Info("Processing request.")

		ctx := newContextWithLogger(req.Context(), logger)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
