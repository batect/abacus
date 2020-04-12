#! /usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
VARS_FILE="$SCRIPT_DIR/../generated.auto.tfvars"

cat <<EOF > "$VARS_FILE"
cloud_sdk_config_name = "$CLOUDSDK_ACTIVE_CONFIG_NAME"
root_domain           = "$ROOT_DOMAIN"
subdomain             = "$SUBDOMAIN"
EOF
