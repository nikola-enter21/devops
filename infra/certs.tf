resource "google_compute_managed_ssl_certificate" "api_cert" {
  name = "gopherify-api-cert"

  managed {
    domains = ["api.users.gopherify.com"]
  }
}
