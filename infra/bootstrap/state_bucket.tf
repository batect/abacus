resource "google_storage_bucket" "state" {
  name               = "${var.project_name}-terraform-state"
  project            = google_project.project.project_id
  location           = var.region
  storage_class      = "REGIONAL"
  bucket_policy_only = true
}
