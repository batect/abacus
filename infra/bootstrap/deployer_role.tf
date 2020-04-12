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
  deployers_group_name = "${google_project.project.name}-deployers@${data.google_organization.organisation.domain}"
}

resource "google_project_iam_custom_role" "deployer" {
  role_id     = "deployer"
  title       = "Deployer"
  description = "Permissions required to deploy the application"
  project     = google_project.project.project_id

  permissions = [
    // Bare minimum required for Terraform to use GCP provider
    "resourcemanager.projects.get",

    // Required to check if Terraform state bucket exists
    "storage.buckets.get",

    // Required to manage GCP project services
    "serviceusage.services.disable",
    "serviceusage.services.enable",
    "serviceusage.services.get",
    "serviceusage.services.list",

    // Required to manage Cloud Run
    "run.services.create",
    "run.services.delete",
    "run.services.get",
    "run.services.getIamPolicy",
    "run.services.setIamPolicy",
    "run.services.update",
    // These permissions are not documented anywhere but are required to manage domain mappings for Cloud Run services.
    "run.domainmappings.create",
    "run.domainmappings.delete",
    "run.domainmappings.get",

    // Required to manage Container Registry storage bucket, and maintain state in Cloud Storage
    "storage.buckets.getIamPolicy",
    "storage.buckets.get",
    "storage.buckets.list",
    "storage.objects.create",
    "storage.objects.delete",
    "storage.objects.get",
    "storage.objects.getIamPolicy",
    "storage.objects.list",
    "storage.objects.setIamPolicy",
    "storage.objects.update",

    // Required to manage Stackdriver uptime checks
    "monitoring.uptimeCheckConfigs.create",
    "monitoring.uptimeCheckConfigs.delete",
    "monitoring.uptimeCheckConfigs.get",
    "monitoring.uptimeCheckConfigs.update",

    // Required to manage Stackdriver notification channels
    "monitoring.notificationChannels.create",
    "monitoring.notificationChannels.delete",
    "monitoring.notificationChannels.get",
    "monitoring.notificationChannels.update",

    // Required to manage Stackdriver alert policies
    "monitoring.alertPolicies.create",
    "monitoring.alertPolicies.delete",
    "monitoring.alertPolicies.get",
    "monitoring.alertPolicies.update",

    // Required to manage BigQuery datasets and tables
    "bigquery.datasets.create",
    "bigquery.datasets.get",
    "bigquery.datasets.getIamPolicy",
    "bigquery.datasets.setIamPolicy",
    "bigquery.datasets.update",
    "bigquery.tables.create",
    "bigquery.tables.get",
    "bigquery.tables.update",
    "bigquery.tables.updateTag",

    // Required to check this IAM role is in sync with configuration
    "iam.roles.get",
    "resourcemanager.projects.getIamPolicy",
    "iam.serviceAccounts.getIamPolicy",
  ]

  depends_on = [google_project_service.iam]
}

resource "google_project_iam_binding" "deployer" {
  role    = google_project_iam_custom_role.deployer.id
  project = google_project.project.project_id
  members = ["group:${local.deployers_group_name}"]
}

