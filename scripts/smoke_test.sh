#! /usr/bin/env bash

set -euo pipefail

function main() {
  read -r -d '' REQUEST_BODY <<EOF || true
  {
    "sessionId": "$(randomUUID)",
    "userId": "$(randomUUID)",
    "sessionStartTime": "$(currentTime)",
    "sessionEndTime": "$(currentTime)",
    "applicationId": "smoke-test-app",
    "applicationVersion": "1.0.0"
  }
EOF

  curl \
    -X PUT \
    -H 'Content-Type: application/json' \
    -d "$REQUEST_BODY" \
    --fail \
    --silent \
    --verbose \
    --show-error \
    "https://$DOMAIN/v1/sessions"
}

function randomUUID() {
  uuidgen | tr '[:upper:]' '[:lower:]' | tr -d '\n'
}

function currentTime() {
  date -u +"%Y-%m-%dT%H:%M:%SZ"
}

main
