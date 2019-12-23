// Copyright 2019 Charles Korn.
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

data "google_dns_managed_zone" "app" {
  name = "app-zone"
}

locals {
  dns_name_without_trailing_dot = trimsuffix(data.google_dns_managed_zone.app.dns_name, ".")
  service_dns_resource_records  = google_cloud_run_domain_mapping.service.status.0.resource_records
}

// Due to https://github.com/terraform-providers/terraform-provider-google/issues/5173, this
// mapping must be manually created and then imported with `terraform import`.
resource "google_cloud_run_domain_mapping" "service" {
  location = google_cloud_run_service.service.location
  name     = "api.${local.dns_name_without_trailing_dot}"

  metadata {
    namespace = google_project_service.cloud_run.project
  }

  spec {
    route_name = google_cloud_run_service.service.name
  }
}

resource "google_dns_record_set" "service" {
  count = length(local.service_dns_resource_records)

  name         = "${google_cloud_run_domain_mapping.service.name}."
  type         = local.service_dns_resource_records[count.index].type
  ttl          = 300
  rrdatas      = [local.service_dns_resource_records[count.index].rrdata]
  managed_zone = data.google_dns_managed_zone.app.name
}
