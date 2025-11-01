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

  gateway_api_config {
    channel = "CHANNEL_STANDARD"
  }

  datapath_provider = "ADVANCED_DATAPATH"

  release_channel {
    channel = "REGULAR"
  }

  monitoring_config {
    enable_components = ["SYSTEM_COMPONENTS", "APISERVER", "SCHEDULER", "CONTROLLER_MANAGER"]
  }

  logging_config {
    enable_components = ["SYSTEM_COMPONENTS", "WORKLOADS"]
  }
}

resource "google_container_node_pool" "primary_nodes" {
  name       = "${var.name_prefix}-nodepool"
  location   = var.zone
  cluster    = google_container_cluster.primary.name
  node_count = 1

  depends_on = [google_container_cluster.primary]

  node_config {
    machine_type    = "e2-medium"
    disk_type       = "pd-balanced"
    disk_size_gb    = 50
    service_account = google_service_account.gke_nodes.email
    oauth_scopes    = ["https://www.googleapis.com/auth/cloud-platform"]

    shielded_instance_config {
      enable_secure_boot          = true
      enable_integrity_monitoring = true
    }

    tags = ["gke-node"]
  }
}

resource "google_service_account" "gke_nodes" {
  account_id   = "${var.name_prefix}-nodes"
  display_name = "GKE Node Service Account"
}
