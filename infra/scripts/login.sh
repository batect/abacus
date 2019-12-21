#! /usr/bin/env bash

set -euo pipefail

LOGIN_TYPE=$1

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CREDS_DIR="$SCRIPT_DIR/../../.creds"
CREDS_PATH="$CREDS_DIR/gcp_service_account_${LOGIN_TYPE}_infra.json"

mkdir -p "$CREDS_DIR"
echo "$GCP_SERVICE_ACCOUNT_KEY" > "$CREDS_PATH"

gcloud auth activate-service-account "$GCP_SERVICE_ACCOUNT_EMAIL" --key-file "$CREDS_PATH"

