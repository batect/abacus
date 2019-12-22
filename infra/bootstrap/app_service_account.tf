resource "google_service_account" "app" {
  account_id   = "batect-abacus-app"
  display_name = "App service account"
  project      = google_project.project.project_id
  depends_on   = [google_project_service.iam]
}

data "google_iam_policy" "app_service_account" {
  binding {
    role = "roles/iam.serviceAccountUser"

    members = ["group:batect-abacus-deployers@googlegroups.com"]
  }
}

resource "google_service_account_iam_policy" "app" {
  service_account_id = google_service_account.app.name
  policy_data        = data.google_iam_policy.app_service_account.policy_data
}
