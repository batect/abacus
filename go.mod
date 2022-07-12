module github.com/batect/abacus

go 1.13

require (
	cloud.google.com/go/storage v1.23.0
	github.com/batect/services-common v0.20.0
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.11.0
	github.com/google/uuid v1.3.0
	github.com/onsi/ginkgo/v2 v2.1.4
	github.com/onsi/gomega v1.19.0
	github.com/sirupsen/logrus v1.8.1
	github.com/unrolled/secure v1.11.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.33.0
	go.opentelemetry.io/otel v1.8.0
	go.opentelemetry.io/otel/trace v1.8.0
	google.golang.org/api v0.87.0
)

// Required until https://github.com/go-playground/validator/pull/601 and https://github.com/go-playground/validator/pull/614 are merged.
replace github.com/go-playground/validator/v10 => github.com/charleskorn/validator/v10 v10.7.1-0.20210711002023-cacc846680e2
