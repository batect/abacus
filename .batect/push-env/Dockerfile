# This Dockerfile is intentionally structured to share as many layers as possible between push-env and terraform.

FROM alpine:3.14.0

RUN apk add --no-cache \
    bash==5.1.0-r0 \
    ca-certificates==20191127-r5 \
    curl==7.77.0-r1 \
    python3==3.8.10-r0

# Based on https://github.com/GoogleCloudPlatform/cloud-sdk-docker/blob/master/alpine/Dockerfile.
ARG CLOUD_SDK_VERSION

RUN curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    tar xzf google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    rm google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz

ENV PATH /google-cloud-sdk/bin:$PATH

RUN gcloud components install beta

RUN apk add --no-cache \
    docker==20.10.3-r1 \
    git==2.30.2-r0
