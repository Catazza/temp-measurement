// ---------------     Initial setup   ---------------   
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

// ---------------   APIs  ---------------   
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

resource "google_project_service" "cloud_build_api" {
  project = var.GCP_PROJECT_ID
  service = "cloudbuild.googleapis.com"

  disable_dependent_services = true
}

// ---------------   GCF terraform ---------------   
data "archive_file" "create_dbloader_zip" {
  type = "zip"
  output_path = "${path.module}/files/compressed_dbloader_func.zip"

  source_dir = "../../functions/dbloader" // TODO: DO BETTER PATH WITH TF FUNCTIONS
}

resource "google_storage_bucket" "staging-gcf-db-loader" {
  name = "staging-gcf-db-loader"
  location = "EU"

  force_destroy = true
}

resource "google_storage_bucket_object" "compressed_dbloader_func" {
  name = format("%s#%s","compressed_dbloader_func", data.archive_file.create_dbloader_zip.output_md5)
  bucket = google_storage_bucket.staging-gcf-db-loader.name
  source = data.archive_file.create_dbloader_zip.output_path 
  content_disposition = "attachment"
  content_encoding    = "gzip"
  content_type        = "application/zip"
}

resource "google_cloudfunctions_function" "dbloader" {
  name = "dbloader"
  description = "function to load sensor measurements from PubSub into Bigquery"
  runtime = "go113"

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.staging-gcf-db-loader.name
  source_archive_object = google_storage_bucket_object.compressed_dbloader_func.name
  timeout               = 60
  entry_point           = "StoreTempMeasurementBQ"

  event_trigger {
    event_type = "google.pubsub.topic.publish"
    resource = "projects/temp-measure-dev/topics/temp-sensor-sink"
    failure_policy {
      retry = "false"
    }
  }
}
// TODO: ADD IAM FOR CLOUD FUNCTION?  

// Second function
data "archive_file" "create_tempreadings_zip" {
  type = "zip"
  output_path = "${path.module}/files/compressed_tempreadings_func.zip"

  source_dir = "../../functions/tempreadings" // TODO: DO BETTER PATH WITH TF FUNCTIONS
}

resource "google_storage_bucket" "staging-gcf-tempreadings" {
  name = "staging-gcf-tempreadings"
  location = "EU"

  force_destroy = true
}

resource "google_storage_bucket_object" "compressed_tempreadings_func" {
  name = format("%s#%s","compressed_tempreadings_func", data.archive_file.create_tempreadings_zip.output_md5)
  bucket = google_storage_bucket.staging-gcf-tempreadings.name
  source = data.archive_file.create_tempreadings_zip.output_path 
  content_disposition = "attachment"
  content_encoding    = "gzip"
  content_type        = "application/zip"
}

resource "google_cloudfunctions_function" "tempreadings" {
  name = "tempreadings"
  description = "function to retrieve last 100 measurements from Bigquery"
  runtime = "go113"

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.staging-gcf-tempreadings.name
  source_archive_object = google_storage_bucket_object.compressed_tempreadings_func.name
  timeout               = 60
  entry_point           = "RetrieveTempreadings"
  trigger_http = true
}

resource "google_cloudfunctions_function_iam_member" "invoker" {
  project        = google_cloudfunctions_function.tempreadings.project
  region         = google_cloudfunctions_function.tempreadings.region
  cloud_function = google_cloudfunctions_function.tempreadings.name

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}

// ---------------   Bigquery storage ---------------   
resource "google_bigquery_dataset" "temp_measure" {
  dataset_id                  = "temp_measure"
  friendly_name               = "temp_measure"
  description                 = "This is the dataset to host the temperature measurements from the raspberry pi"
  location                    = "europe-west2"

  depends_on = [google_project_service.bq_api]
}

resource "google_bigquery_table" "temp_history_raw" {
  dataset_id = google_bigquery_dataset.temp_measure.dataset_id
  table_id = "temp_history_raw"

  time_partitioning {
    type = "DAY"
    field = "processing_time"
  }

  depends_on = [google_bigquery_dataset.temp_measure]
  
  schema = <<EOF
[
  {
    "name": "pubsub_message_id",
    "type": "STRING",
    "description": "the PubSub id of the message"
  },
  {
    "name": "json_msg",
    "type": "STRING",
    "mode": "NULLABLE",
    "description": "The raw json message from the measurement device"
  },
  {
    "name": "processing_time",
    "type": "TIMESTAMP",
    "mode": "NULLABLE",
    "description": "The time the message was processed by the ingestion process"
  }
]
EOF
}

resource "google_bigquery_table" "temp_history_parsed" {
  dataset_id = google_bigquery_dataset.temp_measure.dataset_id
  table_id = "temp_history_parsed"

  time_partitioning {
    type = "DAY"
    field = "processing_time"
  }

  depends_on = [google_bigquery_dataset.temp_measure]
  
  schema = <<EOF
[
  {
    "name": "pubsub_message_id",
    "type": "STRING",
    "description": "the PubSub id of the message"
  },
  {
    "name": "device_message_id",
    "type": "STRING",
    "description": "the id of the message generated by the device"
  },
  {
    "name": "temperature",
    "type": "FLOAT",
    "mode": "NULLABLE",
    "description": "The temperature measurement"
  },
  {
    "name": "humidity",
    "type": "FLOAT",
    "mode": "NULLABLE",
    "description": "The humidity measurement"
  },
  {
    "name": "measurement_time",
    "type": "TIMESTAMP",
    "mode": "NULLABLE",
    "description": "The timestamp the measurement happened on the IOT device"
  },
  {
    "name": "processing_time",
    "type": "TIMESTAMP",
    "mode": "NULLABLE",
    "description": "The time the message was processed by the ingestion process"
  }
]
EOF
}