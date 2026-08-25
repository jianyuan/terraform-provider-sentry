# Retrieve a Jira Server integration
data "sentry_organization_integration" "jira_server" {
  organization = "my-org"
  provider_key = "jira_server"
  name         = "JIRA" # Name of your Jira server
}

resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          jira_server = {
            integration_id = data.sentry_organization_integration.jira_server.id
            project        = "349719"
            issue_type     = "1"

            # Optional Jira fields. All values are IDs, not display names.
            labels     = "oncall,triage" # comma-separated, not a list
            components = ["10001"]
            priority   = "3"
            reporter   = "jira-bot" # Jira Server username

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
