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

locals {
  log_severities_to_ignore         = ["NOTICE", "INFO", "DEBUG", "WARNING"]
  log_severities_for_log_query     = join(" ", formatlist("severity!=\"%s\"", local.log_severities_to_ignore))
  log_severities_for_metrics_query = join(" ", formatlist("metric.label.severity!=\"%s\"", local.log_severities_to_ignore))

  log_names_to_ignore         = ["cloudaudit.googleapis.com/activity", "run.googleapis.com/requests", "monitoring.googleapis.com/ViolationOpenEventv1", "monitoring.googleapis.com/ViolationAutoResolveEventv1"]
  log_names_for_log_query     = join(" ", formatlist("logName!=\"projects/${data.google_project.project.name}/logs/%s\"", [for v in local.log_names_to_ignore : replace(v, "/", "%2F")]))
  log_names_for_metrics_query = join(" ", formatlist("metric.label.log!=\"%s\"", local.log_names_to_ignore))

  log_errors_query = "resource.type=\"cloud_run_revision\" resource.labels.service_name=\"${google_cloud_run_service.service.name}\" ${local.log_severities_for_log_query} ${local.log_names_for_log_query}"
}

resource "google_monitoring_alert_policy" "log_errors" {
  display_name = "Log errors policy"
  combiner     = "OR"

  conditions {
    display_name = "Log entries for service"

    condition_threshold {
      filter          = "metric.type=\"logging.googleapis.com/log_entry_count\" resource.type=\"cloud_run_revision\" resource.label.service_name=\"${google_cloud_run_service.service.name}\" ${local.log_names_for_metrics_query} ${local.log_severities_for_metrics_query}"
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
    **This alert has fired because there was one or more `$${metric.label.severity}` level log messages written to the `$${metric.label.log}` log by `$${resource.label.revision_name}`.**

    [Quick link to logs](https://console.cloud.google.com/logs/query;query=${replace(urlencode(local.log_errors_query), "+", "%20")}?project=${data.google_project.project.name})
    EOT

    mime_type = "text/markdown"
  }
}
