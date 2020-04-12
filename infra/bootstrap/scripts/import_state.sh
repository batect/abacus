#! /usr/bin/env bash

set -euo pipefail

KNOWN_RESOURCES=$([ ! -f terraform.tfstate ] || terraform state list)
HAD_IMPORT_FAILURES=false

read -r -d '' CONTAINER_REGISTRY_HACK_STATE <<EOF || true
{
  "mode": "managed",
  "type": "google_container_registry",
  "name": "registry",
  "provider": "provider.google",
  "instances": [
    {
      "schema_version": 0,
      "attributes": {
        "bucket_self_link": "https://www.googleapis.com/storage/v1/b/artifacts.$GOOGLE_PROJECT.appspot.com",
        "id": "artifacts.$GOOGLE_PROJECT.appspot.com",
        "location": null,
        "project": null
      },
      "private": "bnVsbA=="
    }
  ]
}
EOF

function main() {
  import google_project.project "$GOOGLE_PROJECT"
  import google_project_iam_custom_role.deployer "projects/$GOOGLE_PROJECT/roles/deployer"
  import google_project_iam_custom_role.app_bigquery_access "projects/$GOOGLE_PROJECT/roles/app_bigquery_access"
  import google_project_iam_binding.deployer "$GOOGLE_PROJECT projects/$GOOGLE_PROJECT/roles/deployer"
  import google_storage_bucket.state "$GOOGLE_PROJECT/$GOOGLE_PROJECT-terraform-state"
  import google_project_service.container_registry "$GOOGLE_PROJECT/containerregistry.googleapis.com"
  import google_project_service.iam "$GOOGLE_PROJECT/iam.googleapis.com"
  import google_service_account.service "projects/$GOOGLE_PROJECT/serviceAccounts/service@$GOOGLE_PROJECT.iam.gserviceaccount.com"
  import google_service_account_iam_policy.service "projects/$GOOGLE_PROJECT/serviceAccounts/service@$GOOGLE_PROJECT.iam.gserviceaccount.com"

  # FIXME: Importing a google_container_registry isn't supported (see https://github.com/terraform-providers/terraform-provider-google/issues/6098),
  # so we have to manually add the state to the state file for the time being.
  # When the issue is fixed, use this command:
  # import google_container_registry.registry "$GOOGLE_PROJECT/containerregistry.googleapis.com"
  # Also remove $CONTAINER_REGISTRY_HACK_STATE above and manualImport() below, and remove moreutils from the terraform container
  manualImport google_container_registry.registry "$CONTAINER_REGISTRY_HACK_STATE"

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

function manualImport() {

  if haveStateFor "$1"; then
    echo "Already imported state for $1, skipping."
  else
    echo "Manually adding state for $1..."
    jq ".resources += [$2]" terraform.tfstate | sponge terraform.tfstate || HAD_IMPORT_FAILURES=true
    echo "Done!"
  fi
}

function haveStateFor() {
  contains "$KNOWN_RESOURCES" "$1"
}

function contains() {
  [[ $1 =~ (^|[[:space:]])$2($|[[:space:]]) ]]
}

main
