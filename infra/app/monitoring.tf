resource "google_monitoring_uptime_check_config" "api_ping" {
  display_name = "API (ping)"
  timeout = "10s"
  period = "300s"

  http_check {
    path = "/ping"
    use_ssl = true
    validate_ssl = true
  }

  monitored_resource {
    type = "uptime_url"
    labels = {
      project_id = google_project_service.cloud_run.project
      host = google_cloud_run_domain_mapping.service.name
    }
  }
}
