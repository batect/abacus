resource "google_project_service" "container_registry" {
  service            = "containerregistry.googleapis.com"
  project            = google_project.project.project_id
  disable_on_destroy = false
}

resource "google_project_service" "dns" {
  service            = "dns.googleapis.com"
  project            = google_project.project.project_id
  disable_on_destroy = false
}

resource "google_project_service" "iam" {
  service            = "iam.googleapis.com"
  project            = google_project.project.project_id
  disable_on_destroy = false
}
