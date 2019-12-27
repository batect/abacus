#! /usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
CREDS_FILE="$SCRIPT_DIR/../../../.creds/application_service_account_personal.json"

terraform output application_service_account_key > "$CREDS_FILE"
