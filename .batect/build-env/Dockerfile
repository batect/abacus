FROM golang:1.14.1-stretch

ARG GOLANGCI_LINT_VERSION=1.23.6

RUN cd /usr/local/bin && curl --fail --location --show-error https://github.com/golangci/golangci-lint/releases/download/v$GOLANGCI_LINT_VERSION/golangci-lint-$GOLANGCI_LINT_VERSION-linux-amd64.tar.gz | tar --strip-components=1 --wildcards -xzf - */golangci-lint

ARG batect_cache_setup_command
RUN sh -c "$batect_cache_setup_command"
