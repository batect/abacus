#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
TERRAFORM_DIR="${SCRIPT_DIR}/../tf"
STATE_BUCKET_NAME="${GOOGLE_PROJECT}-terraform-state"

function main {
    echo "Checking project..."
    if ! projectExists; then
        createProject
    fi

    echo "Checking state bucket..."
    if ! stateBucketExists; then
        createStateBucket
    fi

    runTerraformInit
}

function projectExists {
    ! gcloud projects list --filter "$GOOGLE_PROJECT" 2>&1 | grep -q 'Listed 0 items.'
}

function createProject {
    echo "Project does not exist, creating it..."
    gcloud projects create "$GOOGLE_PROJECT"
    gcloud services enable --project "$GOOGLE_PROJECT" cloudresourcemanager.googleapis.com
    gcloud beta billing projects link "$GOOGLE_PROJECT" --billing-account="$GOOGLE_BILLING_ACCOUNT_ID"
}

function stateBucketExists {
    gsutil ls -p "$GOOGLE_PROJECT" -b "gs://$STATE_BUCKET_NAME" >/dev/null 2>&1
}

function createStateBucket {
    echo "State bucket does not exist, creating it..."
    gsutil mb -p "$GOOGLE_PROJECT" -c regional -l "$GOOGLE_REGION" "gs://$STATE_BUCKET_NAME"
    gsutil uniformbucketlevelaccess set on "gs://$STATE_BUCKET_NAME"
}

function runTerraformInit {
    terraform init -input=false -reconfigure -backend-config="bucket=$STATE_BUCKET_NAME" "$TERRAFORM_DIR"
}

main
