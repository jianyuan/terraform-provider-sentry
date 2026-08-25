# Retrieve a Jira integration
data "sentry_organization_integration" "jira" {
  organization = "my-org"
  provider_key = "jira"
  name         = "JIRA" # Name of your Jira server
}

resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          jira = {
            integration_id = data.sentry_organization_integration.jira.id
            project        = "349719"
            issue_type     = "1"

            # Optional Jira fields. All values are IDs, not display names.
            labels     = "oncall,triage" # comma-separated, not a list
            components = ["10001"]
            priority   = "3"
            reporter   = "5b10ac8d82e05b22cc7d4ef5" # Jira account ID

            # Any other Jira field, keyed by field ID. Sentry rewrites
            # camelCase keys on write, so camelCase field IDs must be
            # spelled all-lowercase: `fixversions`, not `fixVersions`.
            additional_fields = {
              customfield_10101 = "sre-team"
              fixversions       = "10500"
            }
          }
        }
      ]
    }
  ]
}
