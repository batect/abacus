resource "google_dns_managed_zone" "app_zone" {
  name     = "app-zone"
  dns_name = "abacus.batect.dev."
  project  = google_project.project.project_id

  dnssec_config {
    state = "on"
  }
}
