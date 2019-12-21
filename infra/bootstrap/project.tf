resource "google_project" "project" {
  name = var.project_name
  project_id = var.project_name
  billing_account = var.billing_account_id

  lifecycle {
    prevent_destroy = true
  }
}
