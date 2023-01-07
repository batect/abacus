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
  bytes_in_kb = 1024
  bytes_in_mb = 1024 * local.bytes_in_kb
  bytes_in_gb = 1024 * local.bytes_in_mb

  logging_free_tier_bytes = 50 * local.bytes_in_gb

  metrics_chargeable_ingestion_free_tier_bytes_per_month  = 150 * local.bytes_in_mb
  metrics_chargeable_ingestion_free_tier_bytes_per_second = local.metrics_chargeable_ingestion_free_tier_bytes_per_month / local.seconds_in_month

  tracing_free_tier_spans = 2500000

  alert_threshold_percentage = 75
  alert_threshold_decimal    = local.alert_threshold_percentage / 100
}

// FIXME: both of these policies are project-scoped but the free tier considers all projects attached to the billing account

resource "google_monitoring_alert_policy" "stackdriver_logging_free_tier" {
  display_name = "Stackdriver Logging free tier"
  combiner     = "OR"

  conditions {
    display_name = "Monthly log bytes ingested"

    condition_threshold {
      filter          = "metric.type=\"logging.googleapis.com/billing/monthly_bytes_ingested\" resource.type=\"global\""
      comparison      = "COMPARISON_GT"
      duration        = "1800s"
      threshold_value = local.logging_free_tier_bytes * local.alert_threshold_decimal

      trigger {
        count = 1
      }

      aggregations {
        alignment_period     = "3600s"
        cross_series_reducer = "REDUCE_SUM"
        per_series_aligner   = "ALIGN_SUM"
      }
    }
  }

  documentation {
    content = "Free tier limit is ${local.logging_free_tier_bytes} bytes per month. This alert fires at ${local.alert_threshold_percentage}% of the free tier threshold. Documentation: https://cloud.google.com/stackdriver/pricing"
  }

  notification_channels = [google_monitoring_notification_channel.email.name]
}

resource "google_monitoring_alert_policy" "stackdriver_tracing_free_tier" {
  display_name = "Stackdriver Tracing free tier"
  combiner     = "OR"

  conditions {
    display_name = "Monthly trace spans ingested"

    condition_threshold {
      filter          = "metric.type=\"cloudtrace.googleapis.com/billing/monthly_spans_ingested\" resource.type=\"global\""
      comparison      = "COMPARISON_GT"
      duration        = "1800s"
      threshold_value = local.tracing_free_tier_spans * local.alert_threshold_decimal

      trigger {
        count = 1
      }

      aggregations {
        alignment_period     = "3600s"
        cross_series_reducer = "REDUCE_SUM"
        per_series_aligner   = "ALIGN_SUM"
      }
    }
  }

  documentation {
    content = "Free tier limit is ${local.tracing_free_tier_spans} spans per month. This alert fires at ${local.alert_threshold_percentage}% of the free tier threshold. Documentation: https://cloud.google.com/stackdriver/pricing"
  }

  notification_channels = [google_monitoring_notification_channel.email.name]
}

resource "google_monitoring_alert_policy" "stackdriver_metrics_free_tier" {
  display_name = "Stackdriver Metrics free tier"
  combiner     = "OR"

  conditions {
    display_name = "Chargeable metrics ingested"

    condition_threshold {
      filter          = "metric.type=\"monitoring.googleapis.com/billing/bytes_ingested\" resource.type=\"global\""
      comparison      = "COMPARISON_GT"
      duration        = "1800s"
      threshold_value = local.metrics_chargeable_ingestion_free_tier_bytes_per_second * local.alert_threshold_decimal

      trigger {
        count = 1
      }

      aggregations {
        alignment_period     = "3600s"
        cross_series_reducer = "REDUCE_SUM"
        per_series_aligner   = "ALIGN_RATE"
      }
    }
  }

  documentation {
    content = "Free tier limit is ${local.metrics_chargeable_ingestion_free_tier_bytes_per_month} bytes per month (${format("%.1f", local.metrics_chargeable_ingestion_free_tier_bytes_per_second)} bytes per second). This alert fires when the ingestion rate exceeds ${local.alert_threshold_percentage}% of the ingestion rate required to exceed the free tier threshold. Documentation: https://cloud.google.com/stackdriver/pricing"
  }

  notification_channels = [google_monitoring_notification_channel.email.name]
}
