# This Dockerfile is intentionally structured to share as many layers as possible between push-env and terraform.

FROM alpine:3.11.2

RUN apk add --no-cache \
    bash \
    ca-certificates \
    curl \
    python2

# Based on https://github.com/GoogleCloudPlatform/cloud-sdk-docker/blob/master/alpine/Dockerfile.
ARG CLOUD_SDK_VERSION

RUN curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    tar xzf google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    rm google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz

ENV PATH /google-cloud-sdk/bin:$PATH

RUN apk add --no-cache \
    docker \
    git

RUN gcloud components install docker-credential-gcr
