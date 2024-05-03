# List a Project's Client Keys
data "sentry_keys" "all" {
  organization = "my-organization"
  project      = "web-app"
}
