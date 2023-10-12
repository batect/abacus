# This Dockerfile is intentionally structured to share as many layers as possible between push-env and terraform.

FROM alpine:3.18.4

RUN apk add --no-cache \
    bash==5.2.15-r5 \
    ca-certificates==20230506-r0 \
    curl==8.3.0-r0 \
    libc6-compat==1.2.4-r1 \
    python3==3.11.5-r0

# Based on https://github.com/GoogleCloudPlatform/cloud-sdk-docker/blob/master/alpine/Dockerfile.
ARG CLOUD_SDK_VERSION

RUN ARCH=$(uname -m | sed 's/aarch64/arm/g') && \
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-$CLOUD_SDK_VERSION-linux-$ARCH.tar.gz && \
    tar xzf google-cloud-sdk-$CLOUD_SDK_VERSION-linux-$ARCH.tar.gz && \
    rm google-cloud-sdk-$CLOUD_SDK_VERSION-linux-$ARCH.tar.gz

ENV PATH /google-cloud-sdk/bin:$PATH

RUN gcloud components install beta

RUN apk add --no-cache \
    docker==23.0.6-r6 \
    git==2.40.1-r0
