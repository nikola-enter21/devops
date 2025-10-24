# VPC Network
resource "google_compute_network" "vpc_network" {
  name                    = "${var.name_prefix}-vpc"
  auto_create_subnetworks = false
}

# Subnetwork with secondary IP ranges for GKE
resource "google_compute_subnetwork" "gke_subnet" {
  name          = "${var.name_prefix}-subnet"
  region        = var.region
  network       = google_compute_network.vpc_network.id
  ip_cidr_range = "10.0.0.0/16"

  secondary_ip_range {
    range_name    = "pods"
    ip_cidr_range = "10.1.0.0/16"
  }

  secondary_ip_range {
    range_name    = "services"
    ip_cidr_range = "10.2.0.0/20"
  }

  private_ip_google_access = true
}
