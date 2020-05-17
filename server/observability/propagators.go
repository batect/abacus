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

package observability

import (
	"context"
	"encoding/binary"
	"regexp"
	"strconv"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
)

type GCPPropagator struct {}

func (g *GCPPropagator) HTTPExtractors() []propagation.HTTPExtractor {
	return []propagation.HTTPExtractor{
		&cloudTraceContextExtractor{},
	}
}

func (g *GCPPropagator) HTTPInjectors() []propagation.HTTPInjector {
	return global.Propagators().HTTPInjectors()
}

type cloudTraceContextExtractor struct {}

func (c *cloudTraceContextExtractor) Extract(ctx context.Context, supplier propagation.HTTPSupplier) context.Context {
	headerValue := supplier.Get("X-Cloud-Trace-Context")
	sc := c.extractSpanContext(headerValue)

	return trace.ContextWithRemoteSpanContext(ctx, sc)
}

// See https://cloud.google.com/trace/docs/setup#force-trace for a description of the X-Cloud-Trace-Context header,
// and https://github.com/census-ecosystem/opencensus-go-exporter-stackdriver/blob/master/propagation/http.go for the OpenCensus implementation.
func (c *cloudTraceContextExtractor) extractSpanContext(headerValue string) trace.SpanContext {
	regex := regexp.MustCompile("^([\\da-fA-F]{32})/(\\d+)(?:;o=([01]))?$")
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
		TraceID: tid,
		SpanID: sidBytes,
		TraceFlags: flags,
	}
}
