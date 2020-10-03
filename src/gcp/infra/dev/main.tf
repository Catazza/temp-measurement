variable "GCP_PROJECT_ID" {
    type = string
    description = "the target project where to run TF"
}

terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
    }
  }

  backend "gcs" {
    bucket  = "tf-state-temp-measure"
    prefix  = "terraform/state"
  }
}

provider "google" {
  version = "3.5.0"

  project = var.GCP_PROJECT_ID
  region  = "europe-west2"
  zone    = "europe-west2-c"
}

resource "google_compute_network" "vpc_network" {
  name = "terraform-network"
}
