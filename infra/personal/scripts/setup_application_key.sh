#! /usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CREDS_DIR="$SCRIPT_DIR/../../../.creds"
APP_CREDS_FILE="$CREDS_DIR/application_service_account_personal.json"
TEST_DRIVER_CREDS_FILE="$CREDS_DIR/test_driver_service_account_personal.json"

terraform output application_service_account_key > "$APP_CREDS_FILE"
terraform output test_driver_service_account_key > "$TEST_DRIVER_CREDS_FILE"
