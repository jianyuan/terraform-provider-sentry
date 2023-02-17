# Retrieve a Metric Alert
# URL format: https://sentry.io/organizations/[organization]/alerts/rules/details/[internal_id]/
data "sentry_metric_alert" "original" {
  organization = "my-organization"
  project      = "my-project"
  internal_id  = "42"
}

# Create a copy of a Metric Alert
resource "sentry_metric_alert" "copy" {
  organization      = data.sentry_metric_alert.original.organization
  project           = data.sentry_metric_alert.original.project
  name              = "${data.sentry_metric_alert.original.name}-copy"
  dataset           = data.sentry_metric_alert.original.dataset
  query             = data.sentry_metric_alert.original.query
  aggregate         = data.sentry_metric_alert.original.aggregate
  time_window       = data.sentry_metric_alert.original.time_window
  threshold_type    = data.sentry_metric_alert.original.threshold_type
  resolve_threshold = data.sentry_metric_alert.original.resolve_threshold

  dynamic "trigger" {
    for_each = data.sentry_metric_alert.original.trigger
    content {
      dynamic "action" {
        for_each = trigger.value.action
        content {
          type              = action.value.type
          target_type       = action.value.target_type
          target_identifier = action.value.target_identifier
          integration_id    = action.value.integration_id
          input_channel_id  = action.value.input_channel_id
        }
      }

      alert_threshold   = trigger.value.alert_threshold
      label             = trigger.value.label
      resolve_threshold = trigger.value.resolve_threshold
      threshold_type    = trigger.value.threshold_type
    }
  }
}
