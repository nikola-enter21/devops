terraform {
  # Specify the required provider and its version
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.8.0"
    }
  }

  # Use remote state to mitigate state corruption in a team environment
  backend "gcs" {
    bucket = "${var.name_prefix}-bucket"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}
