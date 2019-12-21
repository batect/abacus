provider "google" {
  version     = "~> 3.2.0"
  credentials = "${path.module}/../../.creds/gcp_service_account_bootstrap_infra.json"
}

provider "null" {
  version = "~> 2.1"
}
