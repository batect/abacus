// Copyright 2019-2020 Charles Korn.
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

// Note that this account is only used to read data from Cloud Storage - the GCP-internal service account is used to write to BigQuery.
resource "google_service_account" "bigquery_transfer_service" {
  account_id   = "bigquery-transfer-service"
  display_name = "BiqQuery Transfer Service service account for abacus"
  project      = google_project.project.project_id
  depends_on   = [google_project_service.iam]
}

data "google_iam_policy" "bigquery_transfer_service_service_account" {
  binding {
    role = "roles/iam.serviceAccountUser"

    members = ["group:${local.deployers_group_name}"]
  }
}

resource "google_service_account_iam_policy" "bigquery_transfer_service" {
  service_account_id = google_service_account.bigquery_transfer_service.name
  policy_data        = data.google_iam_policy.bigquery_transfer_service_service_account.policy_data
}
