# This is only required until https://github.com/fsouza/fake-gcs-server/pull/300 is merged.
# This is based on https://github.com/fsouza/fake-gcs-server/blob/main/Dockerfile
FROM golang:1.14.6-buster AS builder

WORKDIR /code
RUN git clone --no-checkout -- https://github.com/charleskorn/fake-gcs-server.git . && git -c advice.detachedHead=false checkout 3f67582a13a6dd01ee4fdc7cf4927962e49748cf

ENV CGO_ENABLED=0
RUN go build -o fake-gcs-server

FROM alpine:3.12.0
COPY --from=builder /code/fake-gcs-server /bin/fake-gcs-server
