#! /usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
APP_IMAGE_DIR="$SCRIPT_DIR/../../.batect/app"

HASH=$(git rev-parse HEAD)
IMAGE_NAME="gcr.io/$GOOGLE_PROJECT/abacus"
TAG="$IMAGE_NAME:$HASH"

docker-credential-gcr configure-docker
docker build -t "$TAG" "$APP_IMAGE_DIR"
docker push "$TAG"
