resource "google_project_service" "container_registry" {
  service            = "containerregistry.googleapis.com"
  project            = google_project.project.project_id
  disable_on_destroy = false
}
