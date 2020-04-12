This directory contains two groups of Terraform files:

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

## Creating GCP project for the first time

* Create `batect.local.yml` with your GCP project name and billing account ID:

    ```yaml
    gcpProject: my-project # GCP project ID
    gcpOrganizationId: 123456787890 # GCP organisation ID
    gcpBillingAccountId: 111111-222222-333333 # GCP billing account ID to use
    ```

* Run `./batect setupGCPProject` to create the GCP project
* Run `./batect createGCPBootstrapServiceAccount` to create a service account for use during bootstrapping
* Run `./batect importBootstrapState` to import the project you just created (this will fail as the other resources have not been created yet)
* Run `./batect applyBootstrapTerraform` to create remaining bootstrap resources
* Run `SERVICE_ACCOUNT_NAME=local-deployments ./batect createGCPDeployerServiceAccount` and follow the instructions to save the credentials locally
* Run `SERVICE_ACCOUNT_NAME=github-actions ./batect createGCPDeployerServiceAccount` and follow the instructions to save the credentials on the CI system
