# Retrieve an Issue Alert
data "sentry_issue_alert" "original" {
  organization = "my-organization"
  project      = "my-project"
  id           = "42"
}
