#! /usr/bin/env bash

set -euo pipefail

rm -f infra/bootstrap/terraform-batect-abacus-prod.tfstate
./batect --config-vars-file=batect.prod.yml setupBootstrapTerraform
./batect --config-vars-file=batect.prod.yml importBootstrapState
./batect --config-vars-file=batect.prod.yml planBootstrapTerraform

read -p "Are you sure you want to apply the plan above? (y/N) " -n 1 -r
if [[ $REPLY =~ ^[Yy]$ ]]; then
  ./batect --config-vars-file=batect.prod.yml applyBootstrapTerraform
else
  echo "Cancelled."
  exit 1
fi
