#! /usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

function main() {
  echo "Logging in..."
  gcloud auth login --brief

  echo "Creating service account..."
  gcloud iam service-accounts create "$SERVICE_ACCOUNT_NAME" --project "$GOOGLE_PROJECT"

  echo "Adding service account to deployers group..."
  SERVICE_ACCOUNT_EMAIL="$SERVICE_ACCOUNT_NAME@$GOOGLE_PROJECT.iam.gserviceaccount.com"
  GOOGLE_ORGANIZATION_DOMAIN=$(gcloud organizations describe "$GOOGLE_ORGANIZATION" --format "value(displayName)")
  DEPLOYERS_GROUP_EMAIL="$GOOGLE_PROJECT-deployers@$GOOGLE_ORGANIZATION_DOMAIN"
  echo "FIXME! This is currently disabled due to https://issuetracker.google.com/issues/153767630."
  echo "For the time being, you'll need to add the user $SERVICE_ACCOUNT_EMAIL to the group $DEPLOYERS_GROUP_EMAIL."
  # gcloud beta identity groups memberships add --group-email "$DEPLOYERS_GROUP_EMAIL" --member-email "$SERVICE_ACCOUNT_EMAIL" --organization "$GOOGLE_ORGANIZATION" --project "$GOOGLE_PROJECT"

  echo "Creating access key..."
  SERVICE_ACCOUNT_KEY=$(gcloud iam service-accounts keys create --iam-account "$SERVICE_ACCOUNT_EMAIL" /dev/stdout)

  echo "Logging out..."
  gcloud auth revoke

  echo "How do you want to use these credentials?"
  select answer in "Locally" "On CI"; do
    case $answer in
      "Locally")
        echo "Setting up service account for use with Terraform..."
        CLOUDSDK_ACTIVE_CONFIG_NAME=app-$GOOGLE_PROJECT \
        GCP_SERVICE_ACCOUNT_EMAIL="$SERVICE_ACCOUNT_EMAIL" \
        GCP_SERVICE_ACCOUNT_KEY="$SERVICE_ACCOUNT_KEY" \
          "$SCRIPT_DIR/setup_service_account.sh"

        break
        ;;
      "On CI")
        echo "Configure these environment variables in the CI system:"
        echo "GCP_SERVICE_ACCOUNT_EMAIL=$SERVICE_ACCOUNT_EMAIL"
        echo "GCP_SERVICE_ACCOUNT_KEY=${SERVICE_ACCOUNT_KEY//$'\n'/}"

        break
        ;;
    esac
  done

  echo
  echo "Done."
}

main
