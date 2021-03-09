# Normal workflow

* Setup Terraform with `./batect setupTerraform`
* Plan changes with `./batect planTerraform`
* Apply changes with `./batect applyTerraform`

# Creating GCP project for the first time

* Create `batect.local.yml` with the following details:

    ```yaml
    gcpProject: my-project # GCP project ID
    subdomain: api.abacus.test # Subdomain to for environment, will be <subdomain>.batect.dev (eg. api.abacus.test.batect.dev)
    ```
