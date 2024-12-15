
resource "sentry_issue_alert" "main" {
  organization = sentry_project.main.organization
  project      = sentry_project.main.id
  name         = "My issue alert"

  action_match = "any"
  filter_match = "any"
  frequency    = 30

  conditions_v2 = [
    { first_seen_event = {} },
    { regression_event = {} },
    { reappeared_event = {} },
    { new_high_priority_issue = {} },
    { existing_high_priority_issue = {} },
    {
      event_frequency = {
        comparison_type = "count"
        value           = 100
        interval        = "1h"
      }
    },
    {
      event_frequency = {
        comparison_type     = "percent"
        comparison_interval = "1w"
        value               = 100
        interval            = "1h"
      }
    },
    {
      event_unique_user_frequency = {
        comparison_type = "count"
        value           = 100
        interval        = "1h"
      }
    },
    {
      event_unique_user_frequency = {
        comparison_type     = "percent"
        comparison_interval = "1w"
        value               = 100
        interval            = "1h"
      }
    },
    {
      event_frequency_percent = {
        comparison_type = "count"
        value           = 100
        interval        = "1h"
      }
    },
    {
      event_frequency_percent = {
        comparison_type     = "percent"
        comparison_interval = "1w"
        value               = 100
        interval            = "1h"
      }
    },
  ]

  filters_v2 = [
    {
      age_comparison = {
        comparison_type = "older"
        value           = 10
        time            = "minute"
      }
    },
    {
      issue_occurrences = {
        value = 10
      }
    },
    {
      assigned_to = {
        target_type = "Unassigned"
      }
    },
    {
      assigned_to = {
        target_type       = "Team"
        target_identifier = sentry_team.test.internal_id // Note: This is the internal ID of the team rather than the slug
      }
    },
    {
      latest_adopted_release = {
        oldest_or_newest = "oldest"
        older_or_newer   = "older"
        environment      = "test"
      }
    },
    { latest_release = {} },
    {
      issue_category = {
        value = "Error"
      }
    },
    {
      event_attribute = {
        attribute = "message"
        match     = "CONTAINS"
        value     = "test"
      }
    },
    {
      event_attribute = {
        attribute = "message"
        match     = "IS_SET"
      }
    },
    {
      tagged_event = {
        key   = "key"
        match = "CONTAINS"
        value = "value"
      }
    },
    {
      tagged_event = {
        key   = "key"
        match = "NOT_SET"
      }
    },
    {
      level = {
        match = "EQUAL"
        level = "error"
      }
    },
  ]

  actions_v2 = [/* Please see below for examples */]

}

#
# Send a notification to Suggested Assignees
#

resource "sentry_issue_alert" "member_alert" {
  actions_v2 = [
    {
      notify_email = {
        target_type      = "IssueOwners"
        fallthrough_type = "ActiveMembers"
      }
    },
  ]
  // ...
}

#
# Send a notification to a Member
#

data "sentry_organization_member" "member" {
  organization = data.sentry_organization.test.id
  email        = "test@example.com"
}

resource "sentry_issue_alert" "member_alert" {
  actions_v2 = [
    {
      notify_email = {
        target_type       = "Member"
        target_identifier = data.sentry_organization_member.member.internal_id
        fallthrough_type  = "AllMembers"
      }
    },
  ]
  // ...
}

#
# Send a notification to a Team
#

data "sentry_team" "team" {
  organization = sentry_project.test.organization
  slug         = "my-team"
}

resource "sentry_issue_alert" "team_alert" {
  actions_v2 = [
    {
      notify_email = {
        target_type       = "Team"
        target_identifier = data.sentry_team.team.internal_id
        fallthrough_type  = "AllMembers"
      }
    },
  ]
  // ...
}

#
# Send a Slack notification
#

# Retrieve a Slack integration
data "sentry_organization_integration" "slack" {
  organization = sentry_project.test.organization

  provider_key = "slack"
  name         = "Slack Workspace" # Name of your Slack workspace
}

resource "sentry_issue_alert" "slack_alert" {
  actions_v2 = [
    {
      slack_notify = {
        workspace = data.sentry_organization_integration.slack.id
        channel   = "#warning"
        tags      = "environment,level"
        notes     = "Please <http://example.com|click here> for triage information"
      }
    },
  ]
  // ...
}

#
# Send a Microsoft Teams notification
#

data "sentry_organization_integration" "msteams" {
  organization = sentry_project.test.organization

  provider_key = "msteams"
  name         = "My Team" # Name of your Microsoft Teams team
}

# TODO
resource "sentry_issue_alert" "msteams_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.integrations.msteams.notify_action.MsTeamsNotifyServiceAction",
    "team": ${parseint(data.sentry_organization_integration.msteams.id, 10)},
    "channel": "General"
  }
]
EOT
  // ...
}

#
# Send a Discord notification
#

data "sentry_organization_integration" "discord" {
  organization = sentry_project.test.organization

  provider_key = "discord"
  name         = "Discord Server" # Name of your Discord server
}

# TODO
resource "sentry_issue_alert" "discord_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.integrations.discord.notify_action.DiscordNotifyServiceAction",
    "server": ${parseint(data.sentry_organization_integration.discord.id, 10)},
    "channel_id": 94732897,
    "tags": "browser,user"
  }
]
EOT
  // ...
}

#
# Create a Jira Ticket
#

data "sentry_organization_integration" "jira" {
  organization = sentry_project.test.organization

  provider_key = "jira"
  name         = "JIRA" # Name of your Jira server
}

# TODO
resource "sentry_issue_alert" "jira_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.integrations.jira.notify_action.JiraCreateTicketAction",
    "integration": ${parseint(data.sentry_organization_integration.jira.id, 10)},
    "project": "349719",
    "issueType": "1"
  }
]
EOT
  // ...
}

#
# Create a Jira Server Ticket
#

data "sentry_organization_integration" "jira_server" {
  organization = sentry_project.test.organization

  provider_key = "jira_server"
  name         = "JIRA" # Name of your Jira server
}

# TODO
resource "sentry_issue_alert" "jira_server_alert" {
  actions = <<EOT
[
  {
    "id": "sentry.integrations.jira_server.notify_action.JiraServerCreateTicketAction",
    "integration": ${parseint(data.sentry_organization_integration.jira_server.id, 10)},
    "project": "349719",
    "issueType": "1"
  }
]
EOT
  // ...
}

#
# Create a GitHub Issue
#

data "sentry_organization_integration" "github" {
  organization = sentry_project.test.organization

  provider_key = "github"
  name         = "GitHub"
}

resource "sentry_issue_alert" "github_alert" {
  actions_v2 = [
    {
      github_create_ticket = {
        integration = data.sentry_organization_integration.github.id
        repo        = "default"
        assignee    = "Baxter the Hacker"
        labels      = ["bug", "p1"]
      }
    },
  ]
  // ...
}

#
# Create an Azure DevOps work item
#

data "sentry_organization_integration" "vsts" {
  organization = sentry_project.test.organization

  provider_key = "vsts"
  name         = "Azure DevOps"
}

resource "sentry_issue_alert" "vsts_alert" {
  actions_v2 = [
    {
      azure_devops_create_ticket = {
        integration    = data.sentry_organization_integration.vsts.id
        project        = "0389485"
        work_item_type = "Microsoft.VSTS.WorkItemTypes.Task"
      }
    },
  ]
  // ...
}

#
# Send a PagerDuty notification
#

data "sentry_organization_integration" "pagerduty" {
  organization = sentry_project.test.organization
  provider_key = "pagerduty"
  name         = "PagerDuty"
}

resource "sentry_integration_pagerduty" "pagerduty" {
  organization    = data.sentry_organization_integration.pagerduty.organization
  integration_id  = data.sentry_organization_integration.pagerduty.id
  service         = "issue-alert-service"
  integration_key = "issue-alert-integration-key"
}

resource "sentry_issue_alert" "pagerduty_alert" {
  actions_v2 = [
    {
      pagerduty_notify_service = {
        account  = sentry_integration_pagerduty.pagerduty.integration_id
        service  = sentry_integration_pagerduty.pagerduty.id
        severity = "default"
      }
    },
  ]
  // ...
}

#
# Send an Opsgenie notification
#

data "sentry_organization_integration" "opsgenie" {
  organization = sentry_project.test.organization
  provider_key = "opsgenie"
  name         = "Opsgenie"
}

resource "sentry_integration_opsgenie" "opsgenie" {
  organization    = data.sentry_organization_integration.opsgenie.organization
  integration_id  = data.sentry_organization_integration.opsgenie.id
  team            = "issue-alert-team"
  integration_key = "my-integration-key"
}

resource "sentry_issue_alert" "opsgenie_alert" {
  actions_v2 = [
    {
      opsgenie_notify_team = {
        account  = sentry_integration_opsgenie.opsgenie.integration_id
        team     = sentry_integration_opsgenie.opsgenie.id
        priority = "P1"
      }
    },
  ]
  // ...
}

#
# Send a notification (for all legacy integrations)
#

resource "sentry_issue_alert" "notification_alert" {
  actions_v2 = [
    { notify_event = {} },
  ]
  // ...
}
