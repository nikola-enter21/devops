resource "google_service_account" "terraform_deployer" {
  account_id   = "github-terraform-deployer"
  display_name = "GitHub Actions Terraform Deployer"
}

resource "google_project_iam_member" "terraform_deployer_owner" {
  project = var.project_id
  role    = "roles/owner"
  member  = "serviceAccount:${google_service_account.terraform_deployer.email}"
}

resource "google_iam_workload_identity_pool" "github_pool" {
  workload_identity_pool_id = "github-actions-pool"
  display_name              = "GitHub Actions WIF"
  description               = "OIDC pool for GitHub Actions to impersonate GCP service accounts"
  disabled                  = false
}

resource "google_iam_workload_identity_pool_provider" "github" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-actions"
  display_name                       = "GitHub Actions OIDC Provider"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }

  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.repository" = "assertion.repository"
  }

  attribute_condition = "attribute.repository == \"${var.github_repository}\""
}

resource "google_service_account_iam_member" "terraform_deployer_wif" {
  service_account_id = google_service_account.terraform_deployer.name
  role               = "roles/iam.workloadIdentityUser"

  // principalSet://iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/.../attribute.repository/owner/repo
  member = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github_pool.name}/attribute.repository/${var.github_repository}"
}

output "terraform_deployer_service_account_email" {
  description = "Service account email used by GitHub Actions for Terraform"
  value       = google_service_account.terraform_deployer.email
}

output "github_wif_provider_name" {
  description = "Workload Identity Provider name for GitHub Actions"
  value       = google_iam_workload_identity_pool_provider.github.name
}
