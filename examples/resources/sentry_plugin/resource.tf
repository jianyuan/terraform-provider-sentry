# Create a plugin
resource "sentry_plugin" "default" {
  organization = "my-organization"

  project = "web-app"
  plugin  = "slack"

  config = {
    webhook = "slack://webhook"
  }
}
