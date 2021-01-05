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

data "google_project" "project" {
}

data "google_service_account" "bigquery_transfer_service" {
  account_id = "bigquery-transfer-service"
}

resource "google_bigquery_data_transfer_config" "sessions_transfer" {
  display_name           = "${var.application_id} import"
  data_source_id         = "google_cloud_storage"
  destination_dataset_id = google_bigquery_table.sessions_table.dataset_id
  location               = google_bigquery_table.sessions_table.location
  schedule               = format("every %d hours", var.transfer_job_interval_hours)
  service_account_name   = data.google_service_account.bigquery_transfer_service.email

  params = {
    data_path_template              = "gs://${data.google_project.project.name}-sessions/v1/${var.application_id}/*/*.json"
    destination_table_name_template = google_bigquery_table.sessions_table.table_id
    file_format                     = "JSON"
    max_bad_records                 = 0
    write_disposition               = "APPEND"
  }
}
