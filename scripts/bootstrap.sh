#! /usr/bin/env bash

set -euo pipefail

ENV=$1
VARS_FILE="batect.$ENV.yml"

if [ ! -f "$VARS_FILE" ]; then
  echo "Variables file for environment ($VARS_FILE) does not exist."
  exit 1
fi

PROJECT_NAME=$(yq r "$VARS_FILE" gcpProject)

rm -f "infra/bootstrap/terraform-$PROJECT_NAME.tfstate"
./batect --config-vars-file="$VARS_FILE" setupBootstrapTerraform
./batect --config-vars-file="$VARS_FILE" importBootstrapState
./batect --config-vars-file="$VARS_FILE" planBootstrapTerraform

read -p "Are you sure you want to apply the plan above? (y/N) " -n 1 -r
if [[ $REPLY =~ ^[Yy]$ ]]; then
  ./batect --config-vars-file="$VARS_FILE" applyBootstrapTerraform
else
  echo
  echo "Cancelled."
  exit 1
fi
