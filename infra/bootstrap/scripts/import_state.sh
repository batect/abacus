#! /usr/bin/env bash

set -euo pipefail

STATE_FILE="terraform-${GOOGLE_PROJECT}.tfstate"
KNOWN_RESOURCES=$([ ! -f "$STATE_FILE" ] || terraform state list)
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

read -r -d '' ARTIFACT_REGISTRY_HACK_STATE <<EOF || true
{
  "mode": "managed",
  "type": "google_artifact_registry_repository",
  "name": "images",
  "provider": "provider.google-beta",
  "instances": [
    {
      "schema_version": 0,
      "attributes": {
        "description": "",
        "format": "DOCKER",
        "id": "projects/$GOOGLE_PROJECT/locations/$GOOGLE_REGION/repositories/",
        "kms_key_name": "",
        "labels": {},
        "location": "$GOOGLE_REGION",
        "name": "images",
        "project": "$GOOGLE_PROJECT",
        "repository_id": "images",
        "timeouts": {
          "create": null,
          "delete": null,
          "update": null
        },
      },
      "private": "eyJlMmJmYjczMC1lY2FhLTExZTYtOGY4OC0zNDM2M2JjN2M0YzAiOnsiY3JlYXRlIjoyNDAwMDAwMDAwMDAsImRlbGV0ZSI6MjQwMDAwMDAwMDAwLCJ1cGRhdGUiOjI0MDAwMDAwMDAwMH0sInNjaGVtYV92ZXJzaW9uIjoiMCJ9"
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
  import google_project_service.artifact_registry "$GOOGLE_PROJECT/artifactregistry.googleapis.com"
  import google_project_service.cloud_run "$GOOGLE_PROJECT/run.googleapis.com"
  import google_project_service.container_registry "$GOOGLE_PROJECT/containerregistry.googleapis.com"
  import google_project_service.iam "$GOOGLE_PROJECT/iam.googleapis.com"
  import google_project_service.monitoring "$GOOGLE_PROJECT/monitoring.googleapis.com"
  import google_project_service.profiling "$GOOGLE_PROJECT/cloudprofiler.googleapis.com"
  import google_project_service.stackdriver "$GOOGLE_PROJECT/stackdriver.googleapis.com"
  import google_service_account.service "projects/$GOOGLE_PROJECT/serviceAccounts/service@$GOOGLE_PROJECT.iam.gserviceaccount.com"
  import google_service_account_iam_policy.service "projects/$GOOGLE_PROJECT/serviceAccounts/service@$GOOGLE_PROJECT.iam.gserviceaccount.com"
  import google_project_iam_member.app_bigquery_job_access "$GOOGLE_PROJECT roles/bigquery.jobUser serviceAccount:service@$GOOGLE_PROJECT.iam.gserviceaccount.com"
  import google_project_iam_member.app_tracing_access "$GOOGLE_PROJECT roles/cloudtrace.agent serviceAccount:service@$GOOGLE_PROJECT.iam.gserviceaccount.com"
  import google_project_iam_member.app_profiler_access "$GOOGLE_PROJECT roles/cloudprofiler.agent serviceAccount:service@$GOOGLE_PROJECT.iam.gserviceaccount.com"
  importUsingBetaProvider google_artifact_registry_repository_iam_binding.images_writer "projects/$GOOGLE_PROJECT/locations/$GOOGLE_REGION/repositories/images roles/artifactregistry.writer"

  # FIXME: workaround for https://github.com/terraform-providers/terraform-provider-google/issues/6936
  # When the issue is fixed, used this command:
  # importUsingBetaProvider google_artifact_registry_repository.images "projects/$GOOGLE_PROJECT/locations/$GOOGLE_REGION/repositories/images"
  manualImport google_artifact_registry_repository.images "$ARTIFACT_REGISTRY_HACK_STATE"

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

function importUsingBetaProvider() {
  if haveStateFor "$1"; then
    echo "Already imported state for $1, skipping."
  else
    terraform import -input=false -backup=- -provider=google-beta "$1" "$2" || HAD_IMPORT_FAILURES=true
  fi
}

function manualImport() {
  if haveStateFor "$1"; then
    echo "Already imported state for $1, skipping."
  else
    echo "Manually adding state for $1..."
    jq ".resources += [$2]" "$STATE_FILE" | sponge "$STATE_FILE" || HAD_IMPORT_FAILURES=true
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
