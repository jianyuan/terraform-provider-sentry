# Create a key
resource "sentry_key" "default" {
  organization = "my-organization"

  project = "web-app"
  name    = "My Key"
}
