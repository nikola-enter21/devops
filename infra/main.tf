terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.8.0"
    }
  }

  backend "gcs" {
    bucket = "${var.name_prefix}-bucket"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

resource "google_project_service" "apis" {
  for_each = toset([
    "compute.googleapis.com",              # Needed for VPC, firewall, subnets, etc.
    "container.googleapis.com",            # GKE
    "iam.googleapis.com",                  # IAM roles, service accounts
    "cloudresourcemanager.googleapis.com", # Project-level permissions
    "serviceusage.googleapis.com",         # Enabling/disabling APIs
    "run.googleapis.com",                  # Cloud Run
    "artifactregistry.googleapis.com",     # Artifact Registry
    "cloudbuild.googleapis.com",           # For CI/CD pipelines
    "vpcaccess.googleapis.com",            # VPC Serverless Connector for Cloud Run -> GKE
    "dns.googleapis.com",                  # If you're using custom domain with DNS
  ])

  project            = var.project_id
  service            = each.value
  disable_on_destroy = true
}

