module github.com/batect/abacus

go 1.13

require (
	cloud.google.com/go/storage v1.19.0
	github.com/batect/service-observability v0.13.0
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.10.0
	github.com/google/uuid v1.3.0
	github.com/onsi/ginkgo/v2 v2.1.0
	github.com/onsi/gomega v1.18.0
	github.com/sirupsen/logrus v1.8.1
	github.com/unrolled/secure v1.0.9
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.26.1
	go.opentelemetry.io/otel v1.2.0
	go.opentelemetry.io/otel/trace v1.2.0
	google.golang.org/api v0.66.0
)

// Required until https://github.com/go-playground/validator/pull/601 and https://github.com/go-playground/validator/pull/614 are merged.
replace github.com/go-playground/validator/v10 => github.com/charleskorn/validator/v10 v10.7.1-0.20210711002023-cacc846680e2
