// Copyright 2019-2023 Charles Korn.
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
  schema_template                = file("${path.module}/sessions_schema.json")
  schema_with_session_attributes = replace(local.schema_template, "[\"replaced in Terraform with session attributes\"]", file("${path.module}/${var.table_id}_session_attributes_schema.json"))
  schema_with_event_attributes   = replace(local.schema_with_session_attributes, "[\"replaced in Terraform with event attributes\"]", file("${path.module}/${var.table_id}_event_attributes_schema.json"))
  schema_with_span_attributes    = replace(local.schema_with_event_attributes, "[\"replaced in Terraform with span attributes\"]", file("${path.module}/${var.table_id}_span_attributes_schema.json"))
}

resource "google_bigquery_table" "sessions_table" {
  dataset_id = var.dataset_id
  table_id   = var.table_id

  time_partitioning {
    type                     = "DAY"
    field                    = "sessionStartTime"
    require_partition_filter = true
  }

  schema = local.schema_with_span_attributes

  clustering = [
    "applicationId",
    "applicationVersion",
  ]

  lifecycle {
    prevent_destroy = true
  }
}
