# Retrieve the Default Key
data "sentry_key" "default" {
  organization = "my-organization"

  project = "web-app"
  name    = "Default"
}
