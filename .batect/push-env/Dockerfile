# This Dockerfile is intentionally structured to share as many layers as possible between push-env and terraform.

FROM alpine:3.16.1

RUN apk add --no-cache \
    bash==5.1.16-r2 \
    ca-certificates==20211220-r0 \
    curl==7.83.1-r2 \
    libc6-compat==1.2.3-r0 \
    python3==3.10.5-r0

# Based on https://github.com/GoogleCloudPlatform/cloud-sdk-docker/blob/master/alpine/Dockerfile.
ARG CLOUD_SDK_VERSION

RUN ARCH=$(uname -m | sed 's/aarch64/arm/g') && \
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-$CLOUD_SDK_VERSION-linux-$ARCH.tar.gz && \
    tar xzf google-cloud-sdk-$CLOUD_SDK_VERSION-linux-$ARCH.tar.gz && \
    rm google-cloud-sdk-$CLOUD_SDK_VERSION-linux-$ARCH.tar.gz

ENV PATH /google-cloud-sdk/bin:$PATH

RUN gcloud components install beta

RUN apk add --no-cache \
    docker==20.10.16-r1 \
    git==2.36.2-r0
