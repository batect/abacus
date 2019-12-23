#! /usr/bin/env bash

set -euo pipefail

KNOWN_RESOURCES=$([ ! -f terraform.tfstate ] || terraform state list)

function main() {
  import google_project.project "$GOOGLE_PROJECT"
  import google_project_iam_custom_role.deployer "projects/$GOOGLE_PROJECT/roles/deployer"
  import google_project_iam_binding.deployer "$GOOGLE_PROJECT projects/$GOOGLE_PROJECT/roles/deployer"
  import google_storage_bucket.state "$GOOGLE_PROJECT/$GOOGLE_PROJECT-terraform-state"
  import google_project_service.container_registry "$GOOGLE_PROJECT/containerregistry.googleapis.com"
  import google_project_service.dns "$GOOGLE_PROJECT/dns.googleapis.com"
  import google_project_service.iam "$GOOGLE_PROJECT/iam.googleapis.com"
  import google_service_account.app "projects/$GOOGLE_PROJECT/serviceAccounts/$GOOGLE_PROJECT-app@$GOOGLE_PROJECT.iam.gserviceaccount.com"
  import google_service_account_iam_policy.app "projects/$GOOGLE_PROJECT/serviceAccounts/$GOOGLE_PROJECT-app@$GOOGLE_PROJECT.iam.gserviceaccount.com"
  import google_dns_managed_zone.app_zone app-zone
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
