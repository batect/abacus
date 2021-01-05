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

resource "google_monitoring_uptime_check_config" "api_ping" {
  display_name = "API (ping)"
  timeout      = "10s"
  period       = "300s"

  http_check {
    path         = "/ping"
    use_ssl      = true
    validate_ssl = true
  }

  monitored_resource {
    type = "uptime_url"
    labels = {
      project_id = data.google_project.project.name
      host       = google_cloud_run_domain_mapping.service.name
    }
  }
}

resource "google_monitoring_alert_policy" "api_ping" {
  display_name = "API ping policy"
  combiner     = "OR"

  conditions {
    display_name = "Uptime Health Check on API (ping)"

    condition_threshold {
      filter          = "metric.type=\"monitoring.googleapis.com/uptime_check/check_passed\" resource.type=\"uptime_url\" metric.label.check_id=\"${google_monitoring_uptime_check_config.api_ping.uptime_check_id}\""
      comparison      = "COMPARISON_GT"
      duration        = "300s"
      threshold_value = 1

      trigger {
        count = 1
      }

      aggregations {
        alignment_period     = "600s"
        cross_series_reducer = "REDUCE_COUNT_FALSE"
        group_by_fields      = ["resource.*"]
        per_series_aligner   = "ALIGN_NEXT_OLDER"
      }
    }
  }

  notification_channels = [google_monitoring_notification_channel.email.name]
}
