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

locals {
  one_percent = 0.01
}

resource "google_monitoring_alert_policy" "service_responses" {
  display_name = "API responses"
  combiner     = "OR"

  conditions {
    display_name = "Errors (%)"

    condition_threshold {
      filter             = "metric.type=\"run.googleapis.com/request_count\" resource.type=\"cloud_run_revision\" resource.label.\"service_name\"=\"${google_cloud_run_service.service.name}\" metric.label.\"response_code_class\"!=\"2xx\" metric.label.\"response_code_class\"!=\"3xx\""
      denominator_filter = "metric.type=\"run.googleapis.com/request_count\" resource.type=\"cloud_run_revision\" resource.label.\"service_name\"=\"${google_cloud_run_service.service.name}\""
      comparison         = "COMPARISON_GT"
      duration           = "600s"
      threshold_value    = local.one_percent

      trigger {
        count = 1
      }

      aggregations {
        alignment_period     = "300s"
        cross_series_reducer = "REDUCE_SUM"
        per_series_aligner   = "ALIGN_RATE"
      }

      denominator_aggregations {
        alignment_period     = "300s"
        cross_series_reducer = "REDUCE_SUM"
        per_series_aligner   = "ALIGN_RATE"
      }
    }
  }

  conditions {
    display_name = "Request latency (95th percentile)"

    condition_threshold {
      filter          = "metric.type=\"run.googleapis.com/request_latencies\" resource.type=\"cloud_run_revision\" resource.label.\"service_name\"=\"${google_cloud_run_service.service.name}\""
      comparison      = "COMPARISON_GT"
      duration        = "600s"
      threshold_value = 500

      trigger {
        count = 1
      }

      aggregations {
        alignment_period = "60s"

        // These settings are based on the settings used by the Cloud Run monitoring tab.
        cross_series_reducer = "REDUCE_PERCENTILE_95"
        per_series_aligner   = "ALIGN_DELTA"
      }
    }
  }

  notification_channels = [google_monitoring_notification_channel.email.name]
}
