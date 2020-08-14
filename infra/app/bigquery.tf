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

locals {
  transfer_job_interval_hours = 8
}

data "google_service_account" "bigquery_transfer_service" {
  account_id = "bigquery-transfer-service"
}

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
}

resource "google_bigquery_table" "smoke_test_sessions" {
  dataset_id = google_bigquery_dataset.default.dataset_id
  table_id   = "smoke_test_sessions"

  time_partitioning {
    type                     = "DAY"
    field                    = "sessionStartTime"
    require_partition_filter = true
  }

  schema = replace(file("${path.module}/sessions_schema.json"), "[\"replaced in Terraform with configuration for table\"]", file("${path.module}/smoke_test_sessions_attributes_schema.json"))

  clustering = [
    "applicationId",
    "applicationVersion",
  ]

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_bigquery_data_transfer_config" "smoke_test_transfer" {
  display_name           = "smoke-test-app import"
  data_source_id         = "google_cloud_storage"
  destination_dataset_id = google_bigquery_dataset.default.dataset_id
  location               = google_bigquery_dataset.default.location
  schedule               = format("every %d hours", local.transfer_job_interval_hours)
  service_account_name   = data.google_service_account.bigquery_transfer_service.email

  params = {
    data_path_template              = "gs://${data.google_project.project.name}-sessions/v1/smoke-test-app/*/*.json"
    destination_table_name_template = google_bigquery_table.smoke_test_sessions.table_id
    file_format                     = "JSON"
    max_bad_records                 = 0
    write_disposition               = "APPEND"
  }
}

resource "google_bigquery_table" "batect_sessions" {
  dataset_id = google_bigquery_dataset.default.dataset_id
  table_id   = "batect_sessions"

  time_partitioning {
    type                     = "DAY"
    field                    = "sessionStartTime"
    require_partition_filter = true
  }

  schema = replace(file("${path.module}/sessions_schema.json"), "[\"replaced in Terraform with configuration for table\"]", file("${path.module}/batect_sessions_attributes_schema.json"))

  clustering = [
    "applicationId",
    "applicationVersion",
  ]

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_bigquery_data_transfer_config" "batect_transfer" {
  display_name           = "batect import"
  data_source_id         = "google_cloud_storage"
  destination_dataset_id = google_bigquery_dataset.default.dataset_id
  location               = google_bigquery_dataset.default.location
  schedule               = format("every %d hours", local.transfer_job_interval_hours)
  service_account_name   = data.google_service_account.bigquery_transfer_service.email

  params = {
    data_path_template              = "gs://${data.google_project.project.name}-sessions/v1/batect/*/*.json"
    destination_table_name_template = google_bigquery_table.batect_sessions.table_id
    file_format                     = "JSON"
    max_bad_records                 = 0
    write_disposition               = "APPEND"
  }
}
