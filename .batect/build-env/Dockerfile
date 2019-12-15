FROM golang:1.13.5-stretch

RUN cd /usr/local/bin && curl --fail --location --show-error https://github.com/golangci/golangci-lint/releases/download/v1.21.0/golangci-lint-1.21.0-linux-amd64.tar.gz | tar --strip-components=1 --wildcards -xzf - */golangci-lint
