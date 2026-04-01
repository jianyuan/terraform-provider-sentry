# Send a notification to issue owners. If no issue owners, then send to all members.
resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          email = {
            target_type      = "IssueOwners"
            fallthrough_type = "AllMembers"
          }
        }
      ]
    }
  ]
}

# Send a notification to a team. If no issue owners, then send to all members.
data "sentry_team" "team" { /* ... */ }

resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          email = {
            target_type       = "Team"
            target_identifier = data.sentry_team.team.internal_id
            fallthrough_type  = "AllMembers"
          }
        }
      ]
    }
  ]
}

# Send a notification to a user. If no issue owners, then send to all members.
data "sentry_organization_member" "member" { /* ... */ }

resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          email = {
            target_type       = "Member"
            target_identifier = data.sentry_organization_member.member.internal_id
            fallthrough_type  = "AllMembers"
          }
        }
      ]
    }
  ]
}