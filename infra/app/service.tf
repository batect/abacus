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
