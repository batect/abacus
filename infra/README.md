# Normal workflow

* Configure authentication with `./batect setupGCPBootstrapServiceAccount`
* Setup Terraform with `./batect setupTerraform`
* Plan changes with `./batect planTerraform`
* Apply changes with `./batect applyTerraform`

# Creating GCP project for the first time

* Create `batect.local.yml` with the following details:

    ```yaml
    gcpProject: my-project # GCP project ID
    gcpOrganizationId: 123456787890 # GCP organisation ID
    gcpBillingAccountId: 111111-222222-333333 # GCP billing account ID to use
    subdomain: api.abacus.test # Subdomain to for environment, will be <subdomain>.batect.dev (eg. api.abacus.test.batect.dev)
    ```
