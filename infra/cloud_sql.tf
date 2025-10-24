# Reserve internal IP range for Private Service Access
resource "google_compute_global_address" "sql_psa_range" {
  name          = "${var.name_prefix}-sql-psa-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.main.id
}

# Create the VPC peering connection between your VPC and Google-managed services
resource "google_service_networking_connection" "sql_vpc_connection" {
  network                 = google_compute_network.main.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.sql_psa_range.name]
}

# Cloud SQL PostgreSQL instance (Private IP only)
resource "google_sql_database_instance" "postgres" {
  name             = "${var.name_prefix}-postgres"
  region           = var.region
  database_version = "POSTGRES_15"
  depends_on       = [google_service_networking_connection.sql_vpc_connection]

  settings {
    tier = "db-f1-micro"

    ip_configuration {
      ipv4_enabled                                  = false
      private_network                               = google_compute_network.main.self_link
      enable_private_path_for_google_cloud_services = true
    }

    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
    }

    availability_type           = "ZONAL"
    deletion_protection_enabled = true
  }
}

# Create a default database inside the instance
resource "google_sql_database" "app" {
  name     = "app_db"
  instance = google_sql_database_instance.postgres.name
}

# Create a user for your app to connect with
resource "google_sql_user" "app_user" {
  name     = "appuser"
  instance = google_sql_database_instance.postgres.name
  password = var.db_password
}
