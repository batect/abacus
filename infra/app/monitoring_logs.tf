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
  log_errors_query = "resource.type=\"cloud_run_revision\" resource.labels.service_name=\"${google_cloud_run_service.service.name}\" severity!=\"NOTICE\" severity!=\"ERROR\" severity!=\"WARNING\" logName!=\"projects/${google_project_service.cloud_run.project}/logs/cloudaudit.googleapis.com%2Factivity\" logName!=\"projects/${google_project_service.cloud_run.project}/logs/run.googleapis.com%2Frequests\""
}

resource "google_monitoring_alert_policy" "log_errors" {
  display_name = "Log errors policy"
  combiner     = "OR"

  conditions {
    display_name = "Log entries for service"

    condition_threshold {
      filter          = "metric.type=\"logging.googleapis.com/log_entry_count\" resource.type=\"cloud_run_revision\" resource.label.service_name=\"${google_cloud_run_service.service.name}\" metric.label.log!=\"cloudaudit.googleapis.com/activity\" metric.label.log!=\"run.googleapis.com/requests\" metric.label.severity!=\"NOTICE\" metric.label.severity!=\"INFO\" metric.label.severity!=\"WARNING\""
      comparison      = "COMPARISON_GT"
      threshold_value = 0
      duration        = "0s"

      trigger {
        count = 1
      }

      aggregations {
        alignment_period     = "60s"
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields      = ["metric.label.severity", "resource.label.revision_name", "metric.label.log"]
        per_series_aligner   = "ALIGN_RATE"
      }
    }
  }

  notification_channels = [google_monitoring_notification_channel.email.name]

  documentation {
    content = <<-EOT
    [Quick link to logs](https://console.cloud.google.com/logs/query;query=${replace(urlencode(local.log_errors_query), "+", "%20")}?project=${google_project_service.cloud_run.project})

    Log query details: view logs with the following query at https://console.cloud.google.com/logs/query?project=${google_project_service.cloud_run.project}:

    ```
    ${local.log_errors_query}
    ```
    EOT

    mime_type = "text/markdown"
  }
}
