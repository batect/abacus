// Copyright 2019-2021 Charles Korn.
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

resource "google_bigquery_dataset" "default" {
  dataset_id = "abacus"
  location   = "US"

  access {
    role          = "OWNER"
    special_group = "projectOwners"
  }

  access {
    role          = "WRITER"
    user_by_email = "service-${data.google_project.project.number}@gcp-sa-bigquerydatatransfer.iam.gserviceaccount.com"
  }

  lifecycle {
    prevent_destroy = true
  }
}

module "smoke_test_sessions_table" {
  source         = "./session_table"
  dataset_id     = google_bigquery_dataset.default.dataset_id
  table_id       = "smoke_test_sessions"
  application_id = "smoke-test-app"
}

module "batect_sessions_table" {
  source         = "./session_table"
  dataset_id     = google_bigquery_dataset.default.dataset_id
  table_id       = "batect_sessions"
  application_id = "batect"
}

data "google_service_account" "bigquery_transfer_service" {
  account_id = "bigquery-transfer-service"
}
