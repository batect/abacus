#! /usr/bin/env bash

set -euo pipefail

STATE_FILE="terraform-${GOOGLE_PROJECT}.tfstate"
KNOWN_RESOURCES=$([ ! -f "$STATE_FILE" ] || terraform state list)
HAD_IMPORT_FAILURES=false

function main() {
  import google_project.project "$GOOGLE_PROJECT"
  import google_storage_bucket.bootstrap_state "$GOOGLE_PROJECT/$GOOGLE_PROJECT-bootstrap-terraform-state"

  if [ "$HAD_IMPORT_FAILURES" = "true" ]; then
    echo "One or more resources failed to import, see above."
    exit 1
  fi
}

function import() {
  if haveStateFor "$1"; then
    echo "Already imported state for $1, skipping."
  else
    terraform import -input=false -backup=- "$1" "$2" || HAD_IMPORT_FAILURES=true
  fi
}

function haveStateFor() {
  contains "$KNOWN_RESOURCES" "$1"
}

function contains() {
  [[ $1 =~ (^|[[:space:]])$2($|[[:space:]]) ]]
}

main
