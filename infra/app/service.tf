// Copyright 2019 Charles Korn.
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

resource "google_project_service" "cloud_run" {
  service = "run.googleapis.com"
}

resource "google_cloud_run_service" "service" {
  name     = "abacus"
  location = "us-central1"

  template {
    spec {
      service_account_name = "${google_project_service.cloud_run.project}-app@${google_project_service.cloud_run.project}.iam.gserviceaccount.com"

      containers {
        image = var.image_reference
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [
    google_project_service.cloud_run
  ]
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
