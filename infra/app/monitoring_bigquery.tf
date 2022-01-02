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
  bigquery_transfer_log_query        = "resource.type=\"bigquery_resource\" protoPayload.serviceName=\"bigquery.googleapis.com\" protoPayload.methodName=\"jobservice.jobcompleted\" protoPayload.authenticationInfo.principalEmail:\"gcp-sa-bigquerydatatransfer.iam.gserviceaccount.com\""
  bigquery_transfer_errors_log_query = "${local.bigquery_transfer_log_query} severity=ERROR"

  seconds_in_minute = 60
  seconds_in_hour   = 60 * local.seconds_in_minute
}

resource "google_logging_metric" "bigquery_transfer_jobs" {
  name        = "bigquery_transfer_jobs"
  filter      = local.bigquery_transfer_log_query
  description = "Completed BigQuery transfer jobs. Note that the transfer service only logs to Logging if the job picks up any new files - so this metric will only have data when new files are ingested."

  metric_descriptor {
    display_name = "BigQuery transfer jobs"
    metric_kind  = "DELTA"
    value_type   = "INT64"
    unit         = "1"

    labels {
      key         = "severity"
      value_type  = "STRING"
      description = "Severity of the job completion log message. Info indicates the job succeeded, error indicates the job failed."
    }

    labels {
      key         = "tableId"
      value_type  = "STRING"
      description = "Name of the destination table"
    }

    labels {
      key         = "logName"
      value_type  = "STRING"
      description = "Name of the log this metric was derived from"
    }
  }

  label_extractors = {
    "severity" = "EXTRACT(severity)"
    "tableId"  = "EXTRACT(protoPayload.serviceData.jobCompletedEvent.job.jobConfiguration.load.destinationTable.tableId)"
    "logName"  = "EXTRACT(logName)"
  }
}

resource "google_monitoring_alert_policy" "bigquery_transfer_errors" {
  display_name = "BigQuery transfer jobs"
  combiner     = "OR"

  conditions {
    display_name = "BigQuery transfer job error rate"

    condition_threshold {
      filter             = "metric.type=\"logging.googleapis.com/user/${google_logging_metric.bigquery_transfer_jobs.name}\" resource.type=\"global\" metric.label.severity!=INFO"
      denominator_filter = "metric.type=\"logging.googleapis.com/user/${google_logging_metric.bigquery_transfer_jobs.name}\" resource.type=\"global\""
      comparison         = "COMPARISON_GT"
      duration           = "600s"
      threshold_value    = 0

      trigger {
        count = 1
      }

      aggregations {
        alignment_period     = format("%ds", 2 * module.batect_sessions_table.transfer_job_interval_hours * local.seconds_in_hour)
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields      = ["metric.label.tableId", "metric.label.logName"]
        per_series_aligner   = "ALIGN_RATE"
      }

      denominator_aggregations {
        alignment_period     = format("%ds", 2 * module.batect_sessions_table.transfer_job_interval_hours * local.seconds_in_hour)
        cross_series_reducer = "REDUCE_SUM"
        group_by_fields      = ["metric.label.tableId", "metric.label.logName"]
        per_series_aligner   = "ALIGN_RATE"
      }
    }
  }

  notification_channels = [google_monitoring_notification_channel.email.name]

  documentation {
    content = <<-EOT
    **This alert has fired because there was one or more non-`INFO` level log messages written to the `$${metric.label.logName}` for the BigQuery transfer job to `$${metric.label.tableId}`.**

    [Logs](https://console.cloud.google.com/logs/query;query=${replace(urlencode(local.bigquery_transfer_errors_log_query), "+", "%20")}?project=${data.google_project.project.name}) -
    check the `protoPayload.serviceData.jobCompletedEvent.job.jobStatus` field on the log message for details of the error.

    [Jobs page](https://console.cloud.google.com/bigquery/transfers?project=${data.google_project.project.name})
    EOT

    mime_type = "text/markdown"
  }
}
