// Copyright 2019-2021 Charles Korn.
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

package observability

import (
	"context"
	"encoding/binary"
	"fmt"
	"regexp"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/trace"
)

const headerName = "X-Cloud-Trace-Context"

type GCPPropagator struct{}

func (c GCPPropagator) Extract(ctx context.Context, carrier otel.TextMapCarrier) context.Context {
	headerValue := carrier.Get(headerName)
	sc := c.extractSpanContext(headerValue)

	if sc.IsValid() {
		return trace.ContextWithRemoteSpanContext(ctx, sc)
	}

	return ctx
}

// See https://cloud.google.com/trace/docs/setup#force-trace for a description of the X-Cloud-Trace-Context header,
// and https://github.com/census-ecosystem/opencensus-go-exporter-stackdriver/blob/master/propagation/http.go for the OpenCensus implementation.
func (c GCPPropagator) extractSpanContext(headerValue string) trace.SpanContext {
	regex := regexp.MustCompile(`^([\da-fA-F]{32})/(\d+)(?:;o=([01]))?$`)
	segments := regex.FindStringSubmatch(headerValue)

	if segments == nil {
		return trace.EmptySpanContext()
	}

	tid, err := trace.IDFromHex(segments[1])

	if err != nil {
		return trace.EmptySpanContext()
	}

	sid, err := strconv.ParseUint(segments[2], 10, 64)

	if err != nil {
		return trace.EmptySpanContext()
	}

	sidBytes := trace.SpanID{}
	binary.BigEndian.PutUint64(sidBytes[:], sid)

	flags := byte(0)

	if segments[3] == "1" {
		flags = trace.FlagsSampled
	}

	return trace.SpanContext{
		TraceID:    tid,
		SpanID:     sidBytes,
		TraceFlags: flags,
	}
}

func (c GCPPropagator) Inject(ctx context.Context, carrier otel.TextMapCarrier) {
	sc := trace.SpanFromContext(ctx).SpanContext()
	sid := binary.BigEndian.Uint64(sc.SpanID[:])
	headerValue := fmt.Sprintf("%v/%v", sc.TraceID.String(), sid)

	if sc.IsSampled() {
		headerValue += ";o=1"
	}

	carrier.Set(headerName, headerValue)
}

func (c GCPPropagator) Fields() []string {
	return []string{headerName}
}
