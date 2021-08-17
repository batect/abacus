FROM golang:1.17.0-buster

ARG GOLANGCI_LINT_VERSION=1.41.1

RUN cd /usr/local/bin && curl --fail --location --show-error https://github.com/golangci/golangci-lint/releases/download/v$GOLANGCI_LINT_VERSION/golangci-lint-$GOLANGCI_LINT_VERSION-linux-amd64.tar.gz | tar --strip-components=1 --wildcards -xzf - */golangci-lint
