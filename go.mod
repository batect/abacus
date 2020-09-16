module github.com/batect/abacus

go 1.13

require (
	cloud.google.com/go v0.65.0
	cloud.google.com/go/storage v1.11.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v0.11.0
	github.com/charleskorn/logrus-stackdriver-formatter v0.3.1
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.3.0
	github.com/google/uuid v1.1.2
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/sirupsen/logrus v1.6.0
	github.com/unrolled/secure v1.0.8
	go.opentelemetry.io/contrib/instrumentation/net/http v0.11.0
	go.opentelemetry.io/otel v0.11.0
	go.opentelemetry.io/otel/sdk v0.11.0
	google.golang.org/api v0.32.0
)

// Required until https://github.com/go-playground/validator/pull/601 and https://github.com/go-playground/validator/pull/614 are merged.
replace github.com/go-playground/validator/v10 => github.com/charleskorn/validator/v10 v10.3.1-0.20200523101504-a85cc5797d3d
