// Copyright 2019-2023 Charles Korn.
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

locals {
  cloudflare_zone_id = "b285aeea52df6b888cdee6d2551ebd32" # We can't look this up with a data resource without giving access to all zones in the Cloudflare account :sadface:
  api_dns_fqdn       = "${var.subdomain}.${var.root_domain}"

  # HACK: We only take the first record because Terraform doesn't support dynamic counts (which would be required in the cloudflare_record below)
  service_dns_resource_record = google_cloud_run_domain_mapping.service.status.0.resource_records.0
}

resource "google_cloud_run_domain_mapping" "service" {
  location = google_cloud_run_service.service.location
  name     = local.api_dns_fqdn

  metadata {
    namespace = data.google_project.project.name
  }

  spec {
    route_name = google_cloud_run_service.service.name
  }
}

resource "cloudflare_record" "service" {
  name    = var.subdomain
  type    = local.service_dns_resource_record.type
  value   = trimsuffix(local.service_dns_resource_record.rrdata, ".")
  ttl     = 300
  zone_id = local.cloudflare_zone_id
}
