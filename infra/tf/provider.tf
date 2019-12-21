provider "google" {
  version     = "~> 3.2.0"
  credentials = ".creds/gcp_service_account.json"
}

provider "null" {
  version = "~> 2.1"
}
