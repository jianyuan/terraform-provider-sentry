# Create a plugin
resource "sentry_rule" "default" {
  organization = "my-organization"
  project      = "web-app"

  action_match = "any"
  frequency    = 30
  environment  = "production"

  conditions = [
    {
      id   = "sentry.rules.conditions.first_seen_event.FirstSeenEventCondition"
      name = "A new issue is created"
    }
  ]

  filters = [
    {
      id         = "sentry.rules.filters.assigned_to.AssignedToFilter"
      targetType = "Unassigned"
    }
  ]

  actions = [
    {
      id               = "sentry.mail.actions.NotifyEmailAction"
      name             = "Send an email to IssueOwners"
      targetIdentifier = ""
      targetType       = "IssueOwners"
    }
  ]
}
