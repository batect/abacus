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
  cloud_storage_free_tier_bytes = 5 * local.bytes_in_gb

  cloud_storage_free_tier_class_a_operations_per_month              = 5000
  cloud_storage_free_tier_class_a_operations_per_second             = local.cloud_storage_free_tier_class_a_operations_per_month / local.seconds_in_month
  cloud_run_free_tier_alert_threshold_class_a_operations_per_second = local.cloud_storage_free_tier_class_a_operations_per_second * local.alert_threshold_decimal

  cloud_storage_free_tier_class_b_operations_per_month              = 50000
  cloud_storage_free_tier_class_b_operations_per_second             = local.cloud_storage_free_tier_class_b_operations_per_month / local.seconds_in_month
  cloud_run_free_tier_alert_threshold_class_b_operations_per_second = local.cloud_storage_free_tier_class_b_operations_per_second * local.alert_threshold_decimal

  // FIXME: this is an incomplete list of operations - there are certainly others that aren't included here.
  // Due to the way that filters work in alerts, any operation not listed here will be counted in both alert conditions.
  cloud_storage_class_a_operations           = ["CreateBucket", "ListObjects", "SetIamPolicy", "WriteObject"]
  cloud_storage_class_b_operations           = ["GetBucketMetadata", "GetIamPolicy", "GetObjectMetadata", "ReadObject"]
  cloud_storage_filter_to_class_a_operations = join(" ", formatlist("metric.label.method!=\"%s\"", local.cloud_storage_class_b_operations))
  cloud_storage_filter_to_class_b_operations = join(" ", formatlist("metric.label.method!=\"%s\"", local.cloud_storage_class_a_operations))

  cloud_storage_known_operations = join(", ", sort(concat(local.cloud_storage_class_a_operations, local.cloud_storage_class_b_operations)))

  fifteen_minutes = format("%ds", 15*local.seconds_in_minute)
  six_hours       = format("%ds", 6*local.seconds_in_hour)
}

// FIXME: this policy is project-scoped but the free tier considers all projects attached to the billing account

resource "google_monitoring_alert_policy" "cloud_storage_free_tier" {
  display_name = "Cloud Storage free tier"
  combiner     = "OR"

  conditions {
    display_name = "Storage"

    condition_threshold {
      filter          = "metric.type=\"storage.googleapis.com/storage/total_bytes\" resource.type=\"gcs_bucket\""
      comparison      = "COMPARISON_GT"
      duration        = "1800s"
      threshold_value = local.cloud_storage_free_tier_bytes * local.alert_threshold_decimal

      trigger {
        count = 1
      }

      aggregations {
        alignment_period     = "600s"
        cross_series_reducer = "REDUCE_SUM"
        per_series_aligner   = "ALIGN_SUM"
      }
    }
  }

  conditions {
    display_name = "Class A operations"

    condition_threshold {
      filter          = "metric.type=\"storage.googleapis.com/api/request_count\" resource.type=\"gcs_bucket\" ${local.cloud_storage_filter_to_class_a_operations}"
      comparison      = "COMPARISON_GT"
      duration        = local.six_hours
      threshold_value = local.cloud_run_free_tier_alert_threshold_class_a_operations_per_second

      trigger {
        count = 1
      }

      aggregations {
        alignment_period     = local.fifteen_minutes
        cross_series_reducer = "REDUCE_SUM"
        per_series_aligner   = "ALIGN_RATE"
      }
    }
  }

  conditions {
    display_name = "Class B operations"

    condition_threshold {
      filter          = "metric.type=\"storage.googleapis.com/api/request_count\" resource.type=\"gcs_bucket\" ${local.cloud_storage_filter_to_class_b_operations}"
      comparison      = "COMPARISON_GT"
      duration        = local.six_hours
      threshold_value = local.cloud_run_free_tier_alert_threshold_class_b_operations_per_second

      trigger {
        count = 1
      }

      aggregations {
        alignment_period     = local.fifteen_minutes
        cross_series_reducer = "REDUCE_SUM"
        per_series_aligner   = "ALIGN_RATE"
      }
    }
  }

  documentation {
    content = <<-EOT
    Free tier limits:

    * ${local.cloud_storage_free_tier_bytes} bytes of storage

    * ${local.cloud_storage_free_tier_class_a_operations_per_month} class A requests per calendar month (${format("%.6f", local.cloud_storage_free_tier_class_a_operations_per_second)} requests per second)

      This alert fires when the request rate exceeds ${local.alert_threshold_percentage}% of the request rate required to exceed the free tier threshold.

    * ${local.cloud_storage_free_tier_class_b_operations_per_month} class B requests per calendar month (${format("%.6f", local.cloud_storage_free_tier_class_b_operations_per_second)} requests per second)

      This alert fires when the request rate exceeds ${local.alert_threshold_percentage}% of the request rate required to exceed the free tier threshold.

    Note that this alert will count any operations it doesn't know about (anything other than ${local.cloud_storage_known_operations}) as both a class A and a class B request.

    Documentation: https://cloud.google.com/storage/pricing#cloud-storage-always-free
    EOT

    mime_type = "text/markdown"
  }

  notification_channels = [google_monitoring_notification_channel.email.name]
}

