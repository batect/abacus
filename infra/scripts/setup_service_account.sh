#! /usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CREDS_DIR="$SCRIPT_DIR/../../.creds"
CREDS_PATH="$CREDS_DIR/gcp_service_account_${CLOUDSDK_ACTIVE_CONFIG_NAME}_infra.json"

function main() {
  copyCredsIntoPlace
  createConfiguration
  activateServiceAccount
}

function copyCredsIntoPlace() {
  mkdir -p "$CREDS_DIR"
  echo "$GCP_SERVICE_ACCOUNT_KEY" > "$CREDS_PATH"
}

function createConfiguration() {
  existing=$(gcloud config configurations list --format "value(name)" --filter "name=$CLOUDSDK_ACTIVE_CONFIG_NAME")

  if [[ "$existing" != "$CLOUDSDK_ACTIVE_CONFIG_NAME" ]]; then
    gcloud config configurations create "$CLOUDSDK_ACTIVE_CONFIG_NAME" --no-activate
  fi
}

function activateServiceAccount() {
  gcloud auth activate-service-account "$GCP_SERVICE_ACCOUNT_EMAIL" --key-file "$CREDS_PATH"
}

main
