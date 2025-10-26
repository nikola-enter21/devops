data "google_compute_default_service_account" "default" {
  project = var.project_id
}

output "default_compute_service_account_email" {
  value = data.google_compute_default_service_account.default.email
}
