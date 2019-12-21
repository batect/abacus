resource "google_project_service" "iam" {
  service = "iam.googleapis.com"
}

resource "google_project_iam_custom_role" "deployer" {
  role_id     = "deployer"
  title       = "Deployer"
  description = "Permissions required to deploy the application"
  permissions = [
    // Bare minimum required for Terraform to use GCP provider
    "resourcemanager.projects.get",

    // Required to check if Terraform state bucket exists
    "storage.buckets.get",

    // Required to manage GCP project services
    "serviceusage.services.disable",
    "serviceusage.services.enable",
    "serviceusage.services.get",
    "serviceusage.services.list",

    // Required to manage Cloud Run
    "run.services.create",
    "run.services.delete",
    "run.services.get",
    "run.services.getIamPolicy",
    "run.services.setIamPolicy",
    "run.services.update",

    // Required to manage Container Registry storage bucket, and maintain state in Cloud Storage
    "storage.buckets.getIamPolicy",
    "storage.buckets.setIamPolicy",
    "storage.buckets.get",
    "storage.buckets.list",
    "storage.objects.create",
    "storage.objects.delete",
    "storage.objects.get",
    "storage.objects.getIamPolicy",
    "storage.objects.list",
    "storage.objects.setIamPolicy",
    "storage.objects.update",

    // Required to check this IAM role is in sync with configuration
    "iam.roles.get",
    "resourcemanager.projects.getIamPolicy",
  ]

  depends_on = [google_project_service.iam]
}

resource "google_project_iam_binding" "deployer" {
  role    = google_project_iam_custom_role.deployer.id
  members = ["group:batect-abacus-deployers@googlegroups.com"]
}

