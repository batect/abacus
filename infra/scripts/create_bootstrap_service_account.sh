#! /usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

function main() {
  echo "Logging in..."
  gcloud auth login --brief

  echo "Creating service account..."
  SERVICE_ACCOUNT_NAME="bootstrap"
  gcloud iam service-accounts create "$SERVICE_ACCOUNT_NAME" --project "$GOOGLE_PROJECT"

  echo "Granting service account admin access..."
  SERVICE_ACCOUNT_EMAIL="$SERVICE_ACCOUNT_NAME@$GOOGLE_PROJECT.iam.gserviceaccount.com"
  gcloud projects add-iam-policy-binding "$GOOGLE_PROJECT" --member "serviceAccount:$SERVICE_ACCOUNT_EMAIL" --role roles/owner
  gcloud organizations add-iam-policy-binding "$GOOGLE_ORGANIZATION" --member "serviceAccount:$SERVICE_ACCOUNT_EMAIL" --role roles/resourcemanager.organizationViewer

  echo "Creating access key..."
  SERVICE_ACCOUNT_KEY=$(gcloud iam service-accounts keys create --iam-account "$SERVICE_ACCOUNT_EMAIL" /dev/stdout)

  echo "Logging out..."
  gcloud auth revoke

  echo "Setting up service account for use with Terraform..."
  CLOUDSDK_ACTIVE_CONFIG_NAME=bootstrap \
  GCP_SERVICE_ACCOUNT_EMAIL="$SERVICE_ACCOUNT_EMAIL" \
  GCP_SERVICE_ACCOUNT_KEY="$SERVICE_ACCOUNT_KEY" \
    "$SCRIPT_DIR/setup_service_account.sh"

  echo
  echo "Done."
}

main
