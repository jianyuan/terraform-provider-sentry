resource "sentry_project" "default" {
  organization = "my-organization"

  teams = ["my-first-team", "my-second-team"]
  name  = "web-app"

  platform = "javascript"
}

# Add an App Store Connect source to the project
resource "sentry_project_symbol_source" "http" {
  organization = sentry_project.default.organization
  project      = sentry_project.default.id
  type         = "appStoreConnect"
  name         = "App Store Connect"
  layout = {
    type   = "native"
    casing = "default"
  }

  app_connect_issuer      = "app_connect_issuer"
  app_connect_private_key = <<EOT
-----BEGIN PRIVATE KEY-----
[PRIVATE-KEY]
-----END PRIVATE KEY-----
EOT
  app_id                  = "app_id"
}

# Add a SymbolServer (HTTP) symbol source to the project
resource "sentry_project_symbol_source" "http" {
  organization = sentry_project.default.organization
  project      = sentry_project.default.id
  type         = "http"
  name         = "SymbolServer (HTTP)"
  layout = {
    type   = "native"
    casing = "default"
  }
  url = "https://example.com"
}

# Add a Google Cloud Storage symbol source to the project
resource "sentry_project_symbol_source" "gcs" {
  organization = sentry_project.default.organization
  project      = sentry_project.default.id
  type         = "s3"
  name         = "Google Cloud Storage"
  layout = {
    type   = "native"
    casing = "default"
  }
  bucket       = "gcs-bucket-name"
  client_email = "user@project.iam.gserviceaccount.com"
  private_key  = <<EOT
-----BEGIN PRIVATE KEY-----
[PRIVATE-KEY]
-----END PRIVATE KEY-----
EOT
}

# Add an Amazon S3 symbol source to the project
resource "sentry_project_symbol_source" "s3" {
  organization = sentry_project.default.organization
  project      = sentry_project.default.id
  type         = "s3"
  name         = "Amazon S3"
  layout = {
    type   = "native"
    casing = "default"
  }
  bucket     = "s3-bucket-name"
  region     = "us-east-1"
  access_key = "access_key"
  secret_key = "secret_key"
}
