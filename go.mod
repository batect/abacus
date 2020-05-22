module github.com/batect/abacus

go 1.13

require (
	cloud.google.com/go v0.57.0
	cloud.google.com/go/bigquery v1.8.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v0.1.1-0.20200514210843-966afdc5d38c
	github.com/charleskorn/logrus-stackdriver-formatter v0.3.1
	github.com/go-playground/validator/v10 v10.3.0
	github.com/google/uuid v1.1.1
	github.com/onsi/ginkgo v1.12.2
	github.com/onsi/gomega v1.10.1
	github.com/sirupsen/logrus v1.6.0
	go.opentelemetry.io/otel v0.6.0
	google.golang.org/api v0.25.0
)

replace github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace => github.com/charleskorn/opentelemetry-operations-go/exporter/trace v0.1.1-0.20200517080550-269311d02eaf
