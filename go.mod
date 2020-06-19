module github.com/batect/abacus

go 1.13

require (
	cloud.google.com/go v0.58.0
	cloud.google.com/go/bigquery v1.8.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v0.1.1-0.20200529171817-088b5045d282
	github.com/charleskorn/logrus-stackdriver-formatter v0.3.1
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.3.0
	github.com/google/uuid v1.1.1
	github.com/onsi/ginkgo v1.13.0
	github.com/onsi/gomega v1.10.1
	github.com/sirupsen/logrus v1.6.0
	go.opentelemetry.io/otel v0.6.0
	google.golang.org/api v0.28.0
)

// Required until https://github.com/go-playground/validator/pull/601 and https://github.com/go-playground/validator/pull/614 are merged.
replace github.com/go-playground/validator/v10 => github.com/charleskorn/validator/v10 v10.3.1-0.20200523101504-a85cc5797d3d
