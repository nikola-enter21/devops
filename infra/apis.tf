resource "google_project_service" "apis" {
  for_each = toset([

    # Core GCP
    "compute.googleapis.com",              # Compute Engine
    "container.googleapis.com",            # Google Kubernetes Engine (GKE)
    "iam.googleapis.com",                  # Identity and Access Management
    "cloudresourcemanager.googleapis.com", # Projects, folders, org policies
    "serviceusage.googleapis.com",         # Allows enabling/disabling other APIs

    # CI/CD & Artifacts
    "run.googleapis.com",              # Cloud Run for serverless workloads
    "artifactregistry.googleapis.com", # Container and artifact registry
    "cloudbuild.googleapis.com",       # Cloud Build for CI/CD pipelines

    # Networking & DNS
    "vpcaccess.googleapis.com", # Serverless VPC access connectors
    "dns.googleapis.com",       # Cloud DNS for public/private zones

    # Database
    "sqladmin.googleapis.com", # Cloud SQL Admin API

    # Gateway API + Managed Certificates
    "certificatemanager.googleapis.com", # Cloud Certificate Manager (SSL/TLS certs)
    "networkservices.googleapis.com",    # Gateway API, HTTP(S) LB routing, Traffic Director
    "networksecurity.googleapis.com",    # TLS policies, backend security, mTLS

    # Monitoring & Logging
    "monitoring.googleapis.com", # Cloud Monitoring (metrics)
    "logging.googleapis.com"     # Cloud Logging (logs)
  ])

  project            = var.project_id
  service            = each.value
  disable_on_destroy = false
}
