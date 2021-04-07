# This Dockerfile is intentionally structured to share as many layers as possible between push-env and terraform.

FROM alpine:3.13.4

RUN apk add --no-cache \
    bash==5.1.0-r0 \
    ca-certificates==20191127-r5 \
    curl==7.74.0-r1 \
    python3==3.8.8-r0

# Based on https://github.com/GoogleCloudPlatform/cloud-sdk-docker/blob/master/alpine/Dockerfile.
ARG CLOUD_SDK_VERSION

RUN curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    tar xzf google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    rm google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz

ENV PATH /google-cloud-sdk/bin:$PATH

RUN gcloud components install beta

ENV TERRAFORM_VERSION=0.14.10
RUN wget -q -O /tmp/terraform.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip /tmp/terraform.zip -d /usr/local/bin && \
    rm /tmp/terraform.zip

RUN apk add --no-cache \
    bind-tools==9.16.11-r0 \
    jq==1.6-r1 \
    moreutils==0.65-r0 \
    # For tput:
    ncurses==6.2_p20210109-r0 \
    # For uuidgen:
    util-linux==2.36.1-r1
