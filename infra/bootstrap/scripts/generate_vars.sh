#! /usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
VARS_FILE="$SCRIPT_DIR/../generated.auto.tfvars"

cat <<EOF > "$VARS_FILE"
project_name="$GOOGLE_PROJECT"
billing_account_id="$GOOGLE_BILLING_ACCOUNT_ID"
region="$GOOGLE_REGION"
EOF
