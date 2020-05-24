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
  cloud_run_free_tier_requests_per_month                  = 2000000
  cloud_run_free_tier_alert_threshold_requests_per_month  = local.cloud_run_free_tier_requests_per_month * local.alert_threshold_decimal
  cloud_run_free_tier_alert_threshold_requests_per_second = local.cloud_run_free_tier_alert_threshold_requests_per_month / local.seconds_in_month

  days_in_month    = 31
  seconds_in_day   = 24 * 60 * 60
  seconds_in_month = local.days_in_month * local.seconds_in_day
}

// FIXME: these policies are project-scoped but the free tier considers all projects attached to the billing account

resource "google_monitoring_alert_policy" "cloud_run_free_tier" {
  display_name = "Cloud Run free tier"
  combiner     = "OR"

  conditions {
    display_name = "Requests"

    condition_threshold {
      filter          = "metric.type=\"run.googleapis.com/request_count\" resource.type=\"cloud_run_revision\""
      comparison      = "COMPARISON_GT"
      duration        = "1800s"
      threshold_value = local.cloud_run_free_tier_alert_threshold_requests_per_second

      trigger {
        count = 1
      }

      aggregations {
        alignment_period     = "300s"
        cross_series_reducer = "REDUCE_SUM"
        per_series_aligner   = "ALIGN_RATE"
      }
    }
  }

  documentation {
    content = "Free tier limit is ${local.cloud_run_free_tier_requests_per_month} requests per calendar month. This alert fire when the request rate exceeds ${local.alert_threshold_percentage}% of the request rate required to exceed the free tier threshold. Documentation: https://cloud.google.com/run/pricing"
  }

  notification_channels = [google_monitoring_notification_channel.email.name]
}
