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
  service_name = "abacus"

  # Maximum length of revision name is 63 characters
  service_revision_name = substr("${local.service_name}-${var.image_git_sha}-${regex("@sha256:(.*)$", var.image_reference)[0]}", 0, 63)
}

data "google_service_account" "service" {
  account_id = "service"
}

resource "google_cloud_run_service" "service" {
  name     = local.service_name
  location = "us-central1"

  template {
    spec {
      service_account_name = data.google_service_account.service.email

      containers {
        image = var.image_reference

        env {
          name  = "GOOGLE_PROJECT"
          value = data.google_project.project.name
        }

        env {
          name = "HONEYCOMB_API_KEY"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.honeycomb_api_key.secret_id
              key  = "latest"
            }
          }
        }
      }
    }

    metadata {
      name = local.service_revision_name
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

data "google_iam_policy" "allow_invoke_by_all" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "allow_invoke_by_all" {
  location = google_cloud_run_service.service.location
  service  = google_cloud_run_service.service.name

  policy_data = data.google_iam_policy.allow_invoke_by_all.policy_data
}

resource "google_secret_manager_secret" "honeycomb_api_key" {
  secret_id = "honeycomb-api-key"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_iam_binding" "honeycomb_api_key" {
  secret_id = google_secret_manager_secret.honeycomb_api_key.secret_id
  role      = "roles/secretmanager.secretAccessor"
  members   = ["serviceAccount:${data.google_service_account.service.email}"]
}
