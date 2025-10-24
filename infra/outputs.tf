output "cloud_sql_instance_name" {
  description = "The name of the Cloud SQL instance"
  value       = google_sql_database_instance.postgres.name
}

output "cloud_sql_instance_connection_name" {
  description = "The connection name for Cloud SQL (project:region:instance)"
  value       = google_sql_database_instance.postgres.connection_name
}

output "cloud_sql_private_ip" {
  description = "Private IP address of the Cloud SQL instance"
  value       = google_sql_database_instance.postgres.private_ip_address
}

output "cloud_sql_database_name" {
  description = "The name of the database created in the instance"
  value       = google_sql_database.app.name
}

output "cloud_sql_user_name" {
  description = "The database user created"
  value       = google_sql_user.app_user.name
}

output "default_account" {
  value = data.google_compute_default_service_account.default.email
}
