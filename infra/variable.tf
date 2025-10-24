variable "name_prefix" {
  description = "Prefix for naming all resources"
  type        = string
  default     = "devops-fmi-course"
}

variable "project_id" {
  type        = string
  description = "Your GCP project ID"
}

variable "region" {
  type    = string
  default = "europe-west4"
}

variable "zone" {
  type        = string
  description = "GCP zone (within region)"
  default     = "europe-west4-a"
}
