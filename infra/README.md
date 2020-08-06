This directory contains two groups of Terraform files:

* `bootstrap`: contains Terraform files to create GCP resources required for deployment pipeline to run (eg. creating the Terraform
  state bucket or granting the pipeline user the necessary permissions) or which are too sensitive to allow the GitHub Actions pipeline
  access to perform (eg. IAM role creation). Should only need to be run once, requires GCP admin privileges.
  In the pipeline, a check is performed to ensure that the infrastructure has not drifted from the desired state.

  Configure authentication with `./batect gcpBootstrapLogin` and setup Terraform with `./batect setupBootstrapTerraform`, plan changes with `./batect planBootstrapTerraform`
  and apply changes with `./batect applyBootstrapTerraform`.

* `app`: contains Terraform files to deploy the application. Applied as part of the pipeline.

  Configure authentication with `./batect setupGCPBootstrapServiceAccount`, setup Terraform with `./batect setupTerraform`, plan changes with `./batect planTerraform`
  and apply changes with `./batect applyTerraform`.

## Creating GCP project for the first time

* Create `batect.local.yml` with the following details:

    ```yaml
    gcpProject: my-project # GCP project ID
    gcpOrganizationId: 123456787890 # GCP organisation ID
    gcpBillingAccountId: 111111-222222-333333 # GCP billing account ID to use
    subdomain: api.abacus.test # Subdomain to for environment, will be <subdomain>.batect.dev (eg. api.abacus.test.batect.dev)
    ```

* Run `./batect setupGCPProject` to create the GCP project
* Run `./batect createGCPBootstrapServiceAccount` to create a service account for use during bootstrapping
* Run `./batect setupBootstrapTerraform` to prepare Terraform to bootstrap the new project
* Run `./batect importBootstrapState` to import the project you just created
* Run `./batect applyBootstrapTerraform` to create remaining bootstrap resources
* Run `SERVICE_ACCOUNT_NAME=local-deployments ./batect createGCPDeployerServiceAccount` and follow the instructions to save the credentials locally
* Run `SERVICE_ACCOUNT_NAME=github-actions ./batect createGCPDeployerServiceAccount` and follow the instructions to save the credentials on the CI system

## Switching between projects

You'll need to run the following: `./batect setupBootstrapTerraform`
