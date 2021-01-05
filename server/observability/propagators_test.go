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

package observability_test

import (
	"context"
	"net/http"

	"github.com/batect/abacus/server/observability"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/api/trace/tracetest"
)

// Based on test cases from https://github.com/census-ecosystem/opencensus-go-exporter-stackdriver/blob/master/propagation/http_test.go
var _ = Describe("A GCP tracing propagator", func() {
	propagator := observability.GCPPropagator{}

	Context("when processing incoming requests", func() {
		originalSpanContext := trace.SpanContext{TraceID: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, SpanID: [8]byte{1, 2, 3, 4, 5, 6, 7, 8}}
		originalContext := trace.ContextWithRemoteSpanContext(context.Background(), originalSpanContext)

		Context("given no X-Cloud-Trace-Context header", func() {
			var spanContext trace.SpanContext

			BeforeEach(func() {
				headers := http.Header{}
				ctx := propagator.Extract(originalContext, headers)
				spanContext = trace.RemoteSpanContextFromContext(ctx)
			})

			It("does not modify the span context on the context", func() {
				Expect(spanContext).To(Equal(originalSpanContext))
			})
		})

		Context("given the X-Cloud-Trace-Context header contains a valid trace and span ID", func() {
			var spanContext trace.SpanContext

			BeforeEach(func() {
				headers := http.Header{
					"X-Cloud-Trace-Context": {"105445aa7843bc8bf206b12000100000/18374686479671623803"},
				}

				ctx := propagator.Extract(originalContext, headers)
				spanContext = trace.RemoteSpanContextFromContext(ctx)
			})

			It("returns a span context with the trace and span ID extracted from the header", func() {
				Expect(spanContext).To(Equal(trace.SpanContext{
					TraceID: [16]byte{16, 84, 69, 170, 120, 67, 188, 139, 242, 6, 177, 32, 0, 16, 0, 0},
					SpanID:  [8]byte{255, 0, 0, 0, 0, 0, 0, 123},
				}))
			})
		})

		Context("given the X-Cloud-Trace-Context header contains a valid trace and short span ID", func() {
			var spanContext trace.SpanContext

			BeforeEach(func() {
				headers := http.Header{
					"X-Cloud-Trace-Context": {"105445aa7843bc8bf206b12000100000/123"},
				}

				ctx := propagator.Extract(originalContext, headers)
				spanContext = trace.RemoteSpanContextFromContext(ctx)
			})

			It("returns a span context with the trace and span ID extracted from the header", func() {
				Expect(spanContext).To(Equal(trace.SpanContext{
					TraceID: [16]byte{16, 84, 69, 170, 120, 67, 188, 139, 242, 6, 177, 32, 0, 16, 0, 0},
					SpanID:  [8]byte{0, 0, 0, 0, 0, 0, 0, 123},
				}))
			})
		})

		Context("given the X-Cloud-Trace-Context header contains a valid trace and span ID and explicitly disables tracing", func() {
			var spanContext trace.SpanContext

			BeforeEach(func() {
				headers := http.Header{
					"X-Cloud-Trace-Context": {"105445aa7843bc8bf206b12000100000/18374686479671623803;o=0"},
				}

				ctx := propagator.Extract(originalContext, headers)
				spanContext = trace.RemoteSpanContextFromContext(ctx)
			})

			It("returns a span context with the trace and span ID extracted from the header and no trace flags", func() {
				Expect(spanContext).To(Equal(trace.SpanContext{
					TraceID: [16]byte{16, 84, 69, 170, 120, 67, 188, 139, 242, 6, 177, 32, 0, 16, 0, 0},
					SpanID:  [8]byte{255, 0, 0, 0, 0, 0, 0, 123},
				}))
			})
		})

		Context("given the X-Cloud-Trace-Context header contains a valid trace and span ID and explicitly enables tracing", func() {
			var spanContext trace.SpanContext

			BeforeEach(func() {
				headers := http.Header{
					"X-Cloud-Trace-Context": {"105445aa7843bc8bf206b12000100000/18374686479671623803;o=1"},
				}

				ctx := propagator.Extract(originalContext, headers)
				spanContext = trace.RemoteSpanContextFromContext(ctx)
			})

			It("returns a span context with the trace and span ID extracted from the header and the appropriate trace flag to enable tracing", func() {
				Expect(spanContext).To(Equal(trace.SpanContext{
					TraceID:    [16]byte{16, 84, 69, 170, 120, 67, 188, 139, 242, 6, 177, 32, 0, 16, 0, 0},
					SpanID:     [8]byte{255, 0, 0, 0, 0, 0, 0, 123},
					TraceFlags: trace.FlagsSampled,
				}))
			})
		})

		for _, v := range []string{
			"",
			"/",
			"c1e9153fb27f8ac9f2edac765023676e",
			"c1e9153fb27f8ac9f2edac765023676e/",
			"/13102258660371621412",
			"13102258660371621412",
			"c1e9153fb27f8ac9f2edac765023676e/;",
			"c1e9153fb27f8ac9f2edac765023676e/;o=1",
			"c1e9153fb27f8ac9f2edac765023676e/13102258660371621412;",
			"c1e9153fb27f8ac9f2edac765023676e/13102258660371621412;o",
			"c1e9153fb27f8ac9f2edac765023676e/13102258660371621412;o=",
			"c1e9153fb27f8ac9f2edac765023676e/13102258660371621412;o=2",
		} {
			headerValue := v

			Context("given the X-Cloud-Trace-Context header has the invalid value '"+headerValue+"'", func() {
				var spanContext trace.SpanContext

				BeforeEach(func() {
					headers := http.Header{
						"X-Cloud-Trace-Context": {headerValue},
					}

					ctx := propagator.Extract(originalContext, headers)
					spanContext = trace.RemoteSpanContextFromContext(ctx)
				})

				It("does not modify the span context on the context", func() {
					Expect(spanContext).To(Equal(originalSpanContext))
				})
			})
		}
	})

	Context("when processing outgoing requests", func() {
		var headers http.Header

		BeforeEach(func() {
			headers = http.Header{}

			traceGenerator := func(ctx context.Context) trace.SpanContext {
				return trace.SpanContext{
					TraceID: [16]byte{16, 84, 69, 170, 120, 67, 188, 139, 242, 6, 177, 32, 0, 16, 0, 0},
					SpanID:  [8]byte{255, 0, 0, 0, 0, 0, 0, 123},
				}
			}

			tracer := tracetest.NewTracerProvider(tracetest.WithSpanContextFunc(traceGenerator)).Tracer("Tracer")
			ctx, _ := tracer.Start(context.Background(), "Test trace")
			propagator.Inject(ctx, headers)
		})

		It("adds a X-Cloud-Trace-Context header with the trace ID and span ID", func() {
			Expect(headers).To(HaveKeyWithValue("X-Cloud-Trace-Context", []string{"105445aa7843bc8bf206b12000100000/18374686479671623803"}))
		})
	})
})
