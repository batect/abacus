module github.com/batect/abacus

go 1.13

require (
	cloud.google.com/go/storage v1.14.0
	github.com/batect/service-observability v0.7.0
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.6.0
	github.com/google/uuid v1.2.0
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.12.0
	github.com/sirupsen/logrus v1.8.1
	github.com/unrolled/secure v1.0.8
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.20.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
	google.golang.org/api v0.46.0
)

// Required until https://github.com/go-playground/validator/pull/601 and https://github.com/go-playground/validator/pull/614 are merged.
replace github.com/go-playground/validator/v10 => github.com/charleskorn/validator/v10 v10.3.1-0.20200523101504-a85cc5797d3d
