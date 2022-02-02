// Copyright 2019-2022 Charles Korn.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// and the Commons Clause License Condition v1.0 (the "Condition");
// you may not use this file except in compliance with both the License and Condition.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// You may obtain a copy of the Condition at
//
//     https://commonsclause.com/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License and the Condition is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See both the License and the Condition for the specific language governing permissions and
// limitations under the License and the Condition.

terraform {
  required_providers {
    cloudflare = {
      version = "3.8.0"
      source  = "cloudflare/cloudflare"
    }

    google = {
      version = "4.9.0"
      source  = "hashicorp/google"
    }
  }

  required_version = ">= 0.13"
}

provider "google" {
  credentials = "${path.module}/../../.creds/gcp_service_account_${var.cloud_sdk_config_name}_infra.json"
}

provider "cloudflare" {
  api_token = trimspace(file("${path.module}/../../.creds/cloudflare_key"))
}
