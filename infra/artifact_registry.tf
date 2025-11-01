# Create Docker Artifact Registry
resource "google_artifact_registry_repository" "backend_repo" {
  provider      = google
  location      = var.region
  repository_id = "${var.name_prefix}-repo"
  description   = "Docker repository for backend"
  format        = "DOCKER"
}

# Grant the Artifact Registry Reader role to the default service account.
resource "google_project_iam_member" "gke_artifact_reader" {
  project = var.project_id
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:${google_service_account.gke_nodes.email}"
}
