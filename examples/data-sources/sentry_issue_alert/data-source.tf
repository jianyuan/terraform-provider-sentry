# Retrieve an Issue Alert
# URL format: https://sentry.io/organizations/[organization]/alerts/rules/[project]/[internal_id]/details/
data "sentry_issue_alert" "original" {
  organization = "my-organization"
  project      = "my-project"
  internal_id  = "42"
}

# Create a copy of an Issue Alert
resource "sentry_issue_alert" "copy" {
  organization = data.sentry_issue_alert.original.organization
  project      = data.sentry_issue_alert.original.project

  # Copy and modify attributes as necessary.

  name = "${data.sentry_issue_alert.original.name}-copy"

  action_match = data.sentry_issue_alert.original.action_match
  filter_match = data.sentry_issue_alert.original.filter_match
  frequency    = data.sentry_issue_alert.original.frequency

  conditions = data.sentry_issue_alert.original.conditions
  filters    = data.sentry_issue_alert.original.filters
  actions    = data.sentry_issue_alert.original.actions
}
