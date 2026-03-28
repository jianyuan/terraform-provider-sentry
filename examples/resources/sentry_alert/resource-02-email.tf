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
