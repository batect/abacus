#! /usr/bin/env bash

set -euo pipefail

echo "Waiting for Cloud Run to report as ready..."

while true; do
  conditions=$(gcloud beta run domain-mappings describe --platform managed --domain api.abacus.batect.dev --region "$GOOGLE_REGION" --format json --project "$GOOGLE_PROJECT" | jq '.status.conditions')

  ready=$(echo "$conditions" | jq '.[] | select(.type == "Ready")')
  status=$(echo "$ready" | jq -r '.status')

  if [[ "$status" == "True" ]]; then
    echo "Cloud Run reports as ready!"
    exit 0
  else
    message=$(echo "$ready" | jq -r '.message')
    echo "Status is currently '$status': $message"
    sleep 5
  fi
done

