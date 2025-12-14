variable "name_prefix" {
  description = "Prefix for naming all resources"
  type        = string
}

variable "project_id" {
  type        = string
  description = "Your GCP project ID"
}

variable "region" {
  type = string
}

variable "zone" {
  type        = string
  description = "GCP zone (within region)"
}

variable "db_password" {
  description = "Password for appuser in Postgres"
  type        = string
  sensitive   = true
}

variable "github_repository" {
  description = "GitHub repository in the form owner/name"
  type        = string
}
