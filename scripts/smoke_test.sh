#! /usr/bin/env bash

set -euo pipefail

function main() {
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
    "applicationVersion": "1.0.0"
  }
EOF

  echoBlueText "Sending session..."

  curl \
    -X PUT \
    -H 'Content-Type: application/json' \
    -d "$UPLOAD_DATA" \
    --fail \
    --silent \
    --verbose \
    --show-error \
    "https://$DOMAIN/v1/sessions"

  echo
  echoBlueText "Confirming data was written to BigQuery successfully..."

  RETRIEVED_DATA=$(
    bq \
      --format=json \
      "--project_id=$GOOGLE_PROJECT" \
      query \
      --nouse_legacy_sql \
      "SELECT sessionId, userId, FORMAT_TIMESTAMP(\"%Y-%m-%dT%H:%M:%SZ\", sessionStartTime) AS sessionStartTime, FORMAT_TIMESTAMP(\"%Y-%m-%dT%H:%M:%SZ\", sessionEndTime) AS sessionEndTime, applicationId, applicationVersion FROM $GOOGLE_PROJECT.abacus.sessions WHERE sessionStartTime >= '$CURRENT_DATE' AND sessionId = '$SESSION_ID' AND userID = '$USER_ID';"
  )

  echo

  diff -U 9999 <(echo "[$UPLOAD_DATA]" | jq -S .) <(echo "$RETRIEVED_DATA" | jq -S .) || { echo; echoRedText "Data in BigQuery is not the same as what was submitted. See diff above. '-' represents what was expected, '+' represents what was returned by BigQuery."; exit 1; }

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
