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
  cloud_run_free_tier_requests_per_month                  = 2000000
  cloud_run_free_tier_requests_per_second                 = local.cloud_run_free_tier_requests_per_month / local.seconds_in_month
  cloud_run_free_tier_alert_threshold_requests_per_second = local.cloud_run_free_tier_requests_per_second * local.alert_threshold_decimal

  cloud_run_free_tier_cpu_seconds_per_month                  = 180000
  cloud_run_free_tier_cpu_seconds_per_second                 = local.cloud_run_free_tier_cpu_seconds_per_month / local.seconds_in_month
  cloud_run_free_tier_alert_threshold_cpu_seconds_per_second = local.cloud_run_free_tier_cpu_seconds_per_second * local.alert_threshold_decimal

  cloud_run_free_tier_memory_gb_seconds_per_month                    = 360000
  cloud_run_free_tier_memory_gb_seconds_per_second                   = local.cloud_run_free_tier_memory_gb_seconds_per_month / local.seconds_in_month
  cloud_run_free_tier_alert_threshold_memory_gb_seconds_per_second   = local.cloud_run_free_tier_memory_gb_seconds_per_second * local.alert_threshold_decimal
  cloud_run_free_tier_alert_threshold_memory_byte_seconds_per_second = local.cloud_run_free_tier_alert_threshold_memory_gb_seconds_per_second * local.bytes_in_gb

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

  conditions {
    display_name = "CPU"

    condition_threshold {
      filter          = "metric.type=\"run.googleapis.com/container/billable_instance_time\" resource.type=\"cloud_run_revision\""
      comparison      = "COMPARISON_GT"
      duration        = "1800s"
      threshold_value = local.cloud_run_free_tier_alert_threshold_cpu_seconds_per_second

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

  conditions {
    display_name = "Memory"

    condition_threshold {
      filter          = "metric.type=\"run.googleapis.com/container/memory/allocation_time\" resource.type=\"cloud_run_revision\""
      comparison      = "COMPARISON_GT"
      duration        = "1800s"
      threshold_value = local.cloud_run_free_tier_alert_threshold_memory_byte_seconds_per_second

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
    content = <<-EOT
    Free tier limits:

    * ${local.cloud_run_free_tier_requests_per_month} requests per calendar month (${format("%.6f", local.cloud_run_free_tier_requests_per_second)} requests per second)
    * ${local.cloud_run_free_tier_cpu_seconds_per_month} vCPU seconds per month (${format("%.6f", local.cloud_run_free_tier_cpu_seconds_per_second)} vCPU seconds per second)
    * ${local.cloud_run_free_tier_memory_gb_seconds_per_month} GB seconds per month (${format("%.6f", local.cloud_run_free_tier_memory_gb_seconds_per_second)} GB seconds per second)

    This alert fires when the request rate, vCPU or memory consumption exceeds ${local.alert_threshold_percentage}% of the request rate, vCPU or memory consumption required to exceed the free tier threshold.

    Documentation: https://cloud.google.com/run/pricing
    EOT

    mime_type = "text/markdown"
  }

  notification_channels = [google_monitoring_notification_channel.email.name]
}
