# This Dockerfile is intentionally structured to share as many layers as possible between push-env and terraform.

FROM alpine:3.17.3

RUN apk add --no-cache \
    bash==5.2.15-r0 \
    ca-certificates==20220614-r4 \
    curl==8.0.1-r0 \
    libc6-compat==1.2.3-r4 \
    python3==3.10.11-r0

# Based on https://github.com/GoogleCloudPlatform/cloud-sdk-docker/blob/master/alpine/Dockerfile.
ARG CLOUD_SDK_VERSION

RUN ARCH=$(uname -m | sed 's/aarch64/arm/g') && \
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-$CLOUD_SDK_VERSION-linux-$ARCH.tar.gz && \
    tar xzf google-cloud-sdk-$CLOUD_SDK_VERSION-linux-$ARCH.tar.gz && \
    rm google-cloud-sdk-$CLOUD_SDK_VERSION-linux-$ARCH.tar.gz

ENV PATH /google-cloud-sdk/bin:$PATH

RUN gcloud components install beta

ENV TERRAFORM_VERSION=1.4.6
RUN ARCH=$(uname -m | sed 's/aarch64/arm64/g' | sed 's/x86_64/amd64/g') && \
    wget -q -O /tmp/terraform.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_$ARCH.zip && \
    unzip /tmp/terraform.zip -d /usr/local/bin && \
    rm /tmp/terraform.zip

RUN apk add --no-cache \
    bind-tools==9.18.13-r0 \
    jq==1.6-r2 \
    moreutils==0.67-r0 \
    # For tput:
    ncurses==6.3_p20221119-r0 \
    # For uuidgen:
    util-linux==2.38.1-r1
