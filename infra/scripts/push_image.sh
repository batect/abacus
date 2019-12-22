#! /usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
APP_IMAGE_DIR="$SCRIPT_DIR/../../.batect/app"

HASH=$(git rev-parse HEAD)
IMAGE_NAME="gcr.io/$GOOGLE_PROJECT/abacus"
TAG="$IMAGE_NAME:$HASH"

echo
echo "Configuring Docker credential helper for GCP..."
docker-credential-gcr configure-docker

echo
echo "Building image..."
docker build -t "$TAG" "$APP_IMAGE_DIR"

echo
echo "Pushing image..."
docker push "$TAG"

IMAGE_REFERENCE=$(docker image inspect "$TAG" --format '{{ index .RepoDigests 0 }}')

echo
echo "Image pushed to $IMAGE_REFERENCE, generating Terraform variables file..."

VARS_FILE="$SCRIPT_DIR/../app/image.auto.tfvars"

cat <<EOF > "$VARS_FILE"
image_reference="$IMAGE_REFERENCE"
EOF

echo
echo "Done."
