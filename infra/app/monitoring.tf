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

// If creating this fails with "Error creating NotificationChannel: googleapi: Error 400: 'projects/XXX' is not a workspace.",
// you need to go to the GCP console and open the Monitoring page once to trigger workspace creation.
// This is due to https://github.com/terraform-providers/terraform-provider-google/issues/2605.
resource "google_monitoring_notification_channel" "email" {
  display_name = "Email to alerts@batect.dev"
  type         = "email"

  labels = {
    email_address = "alerts@batect.dev"
  }
}
