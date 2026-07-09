# Send a notification to issue owners. If no issue owners, then send to all members.
resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          email = {
            target_type      = "issue_owners"
            fallthrough_type = "AllMembers"
          }
        }
      ]
    }
  ]
}

# Send a notification to a team.
data "sentry_team" "team" { /* ... */ }

resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          email = {
            target_type = "team"
            target_id   = data.sentry_team.team.internal_id
          }
        }
      ]
    }
  ]
}

# Send a notification to a user.
data "sentry_organization_member" "member" { /* ... */ }

resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          email = {
            target_type = "user"
            target_id   = data.sentry_organization_member.member.internal_id
          }
        }
      ]
    }
  ]
}
