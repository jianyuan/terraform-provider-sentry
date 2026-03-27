resource "sentry_alert" "default" {
  # ...

  action_filters = [
    {
      logic_type = "all"
      actions = [
        {
          pagerduty = {
            integration_id = sentry_integration_pagerduty.pagerduty.integration_id
            service_id     = sentry_integration_pagerduty.pagerduty.id
            service_name   = sentry_integration_pagerduty.pagerduty.service
            severity       = "default"
          }
        }
      ]
    }
  ]
}
