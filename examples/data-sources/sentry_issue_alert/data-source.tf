# Retrieve an Issue Alert
# URL format: https://sentry.io/organizations/[organization]/alerts/rules/[project]/[internal_id]/details/
data "sentry_issue_alert" "original" {
  organization = "my-organization"
  project      = "my-project"
  internal_id  = sentry_issue_alert.test.internal_id
}

# Create a copy of an Issue Alert
resource "sentry_issue_alert" "copy" {
  organization = data.sentry_issue_alert.test.organization
  project      = data.sentry_issue_alert.test.project

  # Copy and modify attributes as necessary.

  name = "${data.sentry_issue_alert.test.name}-copy"

  action_match = data.sentry_issue_alert.test.action_match
  filter_match = data.sentry_issue_alert.test.filter_match
  frequency    = data.sentry_issue_alert.test.frequency

  conditions = data.sentry_issue_alert.test.conditions
  filters    = data.sentry_issue_alert.test.filters
  actions    = data.sentry_issue_alert.test.actions
}
