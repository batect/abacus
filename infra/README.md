This directory contains three groups of Terraform files:

* `bootstrap`: contains Terraform files to create GCP resources required for deployment pipeline to run (eg. creating the Terraform
  state bucket or granting the pipeline user the necessary permissions). Should only need to be run once, requires GCP admin privileges.
  In the pipeline, a check is performed to ensure that the infrastructure has not drifted from the desired state.

  Configure authentication with `./batect gcpBootstrapLogin` and setup Terraform with `./batect setupBootstrapTerraform`.

  Before you can plan or apply changes for the first time, you'll need to import the existing state of the infrastructure with
  `./batect importBootstrapState` (we can't maintain this state remotely as we'd need a bucket to store it in, and this bootstrap step
  is responsible for creating the state bucket - it's a chicken and egg problem).

  Plan changes with `./batect planBootstrapTerraform` and apply changes with `./batect applyBootstrapTerraform`.

* `app`: contains Terraform files to deploy the application. Applied as part of the pipeline.

  Configure authentication with `./batect setupGCPBootstrapServiceAccount`, setup Terraform with `./batect setupTerraform`, plan changes with `./batect planTerraform`
  and apply changes with `./batect applyTerraform`.

* `personal`: contains Terraform files to create an environment that can be used to support local development.

  Configure authentication with `./batect setupGCPPersonalServiceAccount`, setup Terraform with `./batect setupPersonalTerraform`, plan changes with `./batect planPersonalTerraform`
  and apply changes with `./batect applyPersonalTerraform`.
