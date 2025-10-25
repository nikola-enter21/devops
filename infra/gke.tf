

# GKE Cluster with IP Allocation
resource "google_container_cluster" "primary" {
  name     = "${var.name_prefix}-gke"
  location = var.zone

  remove_default_node_pool = true
  initial_node_count       = 1

  network             = google_compute_network.vpc_network.id
  subnetwork          = google_compute_subnetwork.gke_subnet.id
  deletion_protection = false

  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }

  master_auth {
    client_certificate_config {
      issue_client_certificate = false
    }
  }

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }
}

# GKE Node Pool
resource "google_container_node_pool" "primary_nodes" {
  name       = "${var.name_prefix}-nodepool"
  location   = var.zone
  cluster    = google_container_cluster.primary.name
  node_count = 1

  depends_on = [google_container_cluster.primary]

  node_config {
    machine_type = "e2-medium"
    disk_type    = "pd-balanced"
    disk_size_gb = 50
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}
