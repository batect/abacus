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

resource "google_project_service" "cloud_run" {
  service = "run.googleapis.com"
  project = google_project.project.project_id
}

resource "google_project_service" "container_registry" {
  service = "containerregistry.googleapis.com"
  project = google_project.project.project_id
}

resource "google_project_service" "iam" {
  service            = "iam.googleapis.com"
  project            = google_project.project.project_id
  disable_on_destroy = false
}

resource "google_project_service" "monitoring" {
  service = "monitoring.googleapis.com"
  project = google_project.project.project_id
}

resource "google_project_service" "profiling" {
  service = "cloudprofiler.googleapis.com"
  project = google_project.project.project_id
}

resource "google_project_service" "stackdriver" {
  service = "stackdriver.googleapis.com"
  project = google_project.project.project_id
}
