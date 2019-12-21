# See https://github.com/terraform-providers/terraform-provider-google/issues/4184 for an explanation of this.

resource "google_project_service" "container_registry" {
  service = "containerregistry.googleapis.com"
}

resource "null_resource" "initialise_container_registry" {
  provisioner "local-exec" {
    command = <<EOF
      docker-credential-gcr configure-docker && \
      (echo 'FROM scratch'; echo 'LABEL maintainer=charleskorn.com') | docker build -t gcr.io/${google_project_service.container_registry.project}/scratch:latest - && \
      docker push gcr.io/${google_project_service.container_registry.project}/scratch:latest
EOF
  }

  depends_on = [google_project_service.container_registry]
}
