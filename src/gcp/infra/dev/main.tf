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

resource "google_project_service" "iot_api" {
  project = var.GCP_PROJECT_ID
  service = "cloudiot.googleapis.com"

  disable_dependent_services = true
}

resource "google_project_service" "bq_api" {
  project = var.GCP_PROJECT_ID
  service = "bigquery.googleapis.com"

  disable_dependent_services = true
}

resource "google_bigquery_dataset" "temp_measure" {
  dataset_id                  = "temp_measure"
  friendly_name               = "temp_measure"
  description                 = "This is the dataset to host the temperature measurements from the raspberry pi"
  location                    = "EU"

  depends_on = [google_project_service.bq_api]
}
