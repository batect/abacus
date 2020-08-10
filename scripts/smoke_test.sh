#! /usr/bin/env bash

set -euo pipefail

BASE_URL=${1:-https://$DOMAIN}

function main() {
  echoBlueText "Generating data..."

  CURRENT_DATE=$(currentDate)
  SESSION_ID=$(randomUUID)
  USER_ID=$(randomUUID)

  read -r -d '' UPLOAD_DATA <<EOF || true
  {
    "sessionId": "$SESSION_ID",
    "userId": "$USER_ID",
    "sessionStartTime": "$(currentTime)",
    "sessionEndTime": "$(currentTime)",
    "applicationId": "smoke-test-app",
    "applicationVersion": "1.0.0",
    "attributes": {
      "source": "smoke-test"
    }
  }
EOF

  echo "Generated session:"
  echo "$UPLOAD_DATA"
  echo

  echoBlueText "Sending session..."

  curl \
    -X PUT \
    -H 'Content-Type: application/json' \
    -d "$UPLOAD_DATA" \
    --fail \
    --silent \
    --verbose \
    --show-error \
    "$BASE_URL/v1/sessions"

  echo
  echoBlueText "Confirming data was written to Cloud Storage successfully..."

  RETRIEVED_DATA=$(gsutil cat "gs://$GOOGLE_PROJECT-sessions/v1/smoke-test-app/1.0.0/$SESSION_ID.json")

  echo
  echo "Response from Cloud Storage: "
  echo "$RETRIEVED_DATA"
  echo

  diff -U 9999 <(echo "$UPLOAD_DATA" | jq -S .) <(echo "$RETRIEVED_DATA" | jq -S 'del(.ingestionTime)') || { echo; echoRedText "Data in Cloud Storage is not the same as what was submitted. See diff above. '-' represents what was expected, '+' represents what was returned by Cloud Storage."; exit 1; }

  echoGreenText "Smoke test completed successfully."
}

function randomUUID() {
  uuidgen | tr '[:upper:]' '[:lower:]' | tr -d '\n'
}

function currentDate() {
  date -u +"%Y-%m-%d"
}

function currentTime() {
  date -u +"%Y-%m-%dT%H:%M:%SZ"
}

function echoBlueText() {
  echo "$(tput setaf 4)$1$(tput sgr0)"
}

function echoGreenText() {
  echo "$(tput setaf 2)$1$(tput sgr0)"
}

function echoRedText() {
  echo "$(tput setaf 1)$1$(tput sgr0)"
}

main
