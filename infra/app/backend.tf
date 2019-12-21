terraform {
  backend "gcs" {
    // This backend will be configured automatically by initialize.sh.
    credentials = "${path.module}/../../.creds/gcp_service_account_app_infra.json"
  }
}
