terraform {
  backend "gcs" {
    // This backend will be configured automatically by initialize.sh.
    credentials = ".creds/gcp_service_account.json"
  }
}
