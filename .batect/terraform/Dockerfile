# This Dockerfile is intentionally structured to share as many layers as possible between push-env and terraform.

FROM alpine:3.16.2

RUN apk add --no-cache \
    bash==5.1.16-r2 \
    ca-certificates==20220614-r0 \
    curl==7.83.1-r4 \
    libc6-compat==1.2.3-r1 \
    python3==3.10.8-r0

# Based on https://github.com/GoogleCloudPlatform/cloud-sdk-docker/blob/master/alpine/Dockerfile.
ARG CLOUD_SDK_VERSION

RUN ARCH=$(uname -m | sed 's/aarch64/arm/g') && \
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-$CLOUD_SDK_VERSION-linux-$ARCH.tar.gz && \
    tar xzf google-cloud-sdk-$CLOUD_SDK_VERSION-linux-$ARCH.tar.gz && \
    rm google-cloud-sdk-$CLOUD_SDK_VERSION-linux-$ARCH.tar.gz

ENV PATH /google-cloud-sdk/bin:$PATH

RUN gcloud components install beta

ENV TERRAFORM_VERSION=1.3.4
RUN ARCH=$(uname -m | sed 's/aarch64/arm64/g' | sed 's/x86_64/amd64/g') && \
    wget -q -O /tmp/terraform.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_$ARCH.zip && \
    unzip /tmp/terraform.zip -d /usr/local/bin && \
    rm /tmp/terraform.zip

RUN apk add --no-cache \
    bind-tools==9.16.33-r0 \
    jq==1.6-r1 \
    moreutils==0.67-r0 \
    # For tput:
    ncurses==6.3_p20220521-r0 \
    # For uuidgen:
    util-linux==2.38-r1
