#! /usr/bin/env bash

set -euo pipefail

KNOWN_RESOURCES=$([ ! -f terraform.tfstate ] || terraform state list)

function main() {
  import google_project.project "$GOOGLE_PROJECT"
  import google_project_iam_custom_role.deployer "projects/$GOOGLE_PROJECT/roles/deployer"
  import google_project_iam_binding.deployer "$GOOGLE_PROJECT roles/deployer"
  import google_storage_bucket.state "$GOOGLE_PROJECT/batect-abacus-terraform-state"
  import google_project_service.container_registry "$GOOGLE_PROJECT/containerregistry.googleapis.com"
  import google_project_service.iam "$GOOGLE_PROJECT/iam.googleapis.com"
}

function import() {
  if haveStateFor "$1"; then
    echo "Already imported state for $1, skipping."
  else
    terraform import -input=false -backup=- "$1" "$2"
  fi
}

function haveStateFor() {
  contains "$KNOWN_RESOURCES" "$1"
}

function contains() {
  [[ $1 =~ (^|[[:space:]])$2($|[[:space:]]) ]]
}

main
