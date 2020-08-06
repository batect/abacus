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

resource "google_bigquery_dataset" "default" {
  dataset_id = "abacus"
  location   = "US"

  access {
    role          = "OWNER"
    special_group = "projectOwners"
  }
}

resource "google_bigquery_table" "sessions" {
  dataset_id = google_bigquery_dataset.default.dataset_id
  table_id   = "sessions"

  time_partitioning {
    type                     = "DAY"
    field                    = "sessionStartTime"
    require_partition_filter = true
  }

  schema = file("${path.module}/sessions_schema.json")

  clustering = [
    "applicationId",
    "applicationVersion",
  ]

  lifecycle {
    prevent_destroy = true
  }
}
