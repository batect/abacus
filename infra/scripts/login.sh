#! /usr/bin/env bash

set -euo pipefail

gcloud auth activate-service-account "$GCP_SERVICE_ACCOUNT_EMAIL" --key-file <( echo "$GCP_SERVICE_ACCOUNT_KEY" )
