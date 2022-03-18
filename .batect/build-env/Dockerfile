FROM golang:1.18.0-buster

ARG GOLANGCI_LINT_VERSION=1.45.0

RUN cd /usr/local/bin && curl --fail --location --show-error https://github.com/golangci/golangci-lint/releases/download/v$GOLANGCI_LINT_VERSION/golangci-lint-$GOLANGCI_LINT_VERSION-linux-$(uname -m | sed 's/aarch64/arm64/g' | sed 's/x86_64/amd64/g' ).tar.gz | tar --strip-components=1 --wildcards -xzf - */golangci-lint
