# Create Docker Artifact Registry
resource "google_artifact_registry_repository" "backend_repo" {
  provider      = google
  location      = var.region
  repository_id = "${var.name_prefix}-repo"
  description   = "Docker repository for backend"
  format        = "DOCKER"
}

# IAM permission for GKE to pull images
resource "google_project_iam_member" "gke_artifact_reader" {
  project = var.project_id
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:service-${data.google_project.project.number}@compute.googleapis.com"
}
