// Copyright 2019-2020 Charles Korn.
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

resource "google_container_registry" "registry" {
  // Nothing to configure.

  depends_on = [google_project_service.container_registry]
}

resource "google_artifact_registry_repository" "images" {
  provider = google-beta

  location      = var.region
  repository_id = "images"
  format        = "DOCKER"

  depends_on = [google_project_service.artifact_registry]
}

resource "google_artifact_registry_repository_iam_binding" "images_writer" {
  provider = google-beta

  location   = google_artifact_registry_repository.images.location
  repository = google_artifact_registry_repository.images.name
  role       = "roles/artifactregistry.writer"
  members    = ["group:${local.deployers_group_name}"]
}
