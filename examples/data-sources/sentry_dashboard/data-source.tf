# Retrieve a Dashboard
# URL format: https://sentry.io/organizations/[organization]/dashboard/[internal_id]/
data "sentry_dashboard" "original" {
  organization = "my-organization"
  internal_id  = "42"
}

# Create a copy of a Dashboard
resource "sentry_dashboard" "copy" {
  organization = data.sentry_dashboard.original.organization

  # Copy and modify attributes as necessary.

  title = "${data.sentry_dashboard.original.title}-copy"

  dynamic "widget" {
    for_each = data.sentry_dashboard.original.widget
    content {
      title        = widget.value.title
      display_type = widget.value.display_type
      interval     = widget.value.interval
      widget_type  = widget.value.widget_type
      limit        = widget.value.limit

      dynamic "query" {
        for_each = widget.value.query
        content {
          name = query.value.name

          fields        = query.value.fields
          aggregates    = query.value.aggregates
          columns       = query.value.columns
          field_aliases = query.value.field_aliases
          conditions    = query.value.conditions
          order_by      = query.value.order_by
        }
      }

      layout {
        x     = widget.value.layout[0].x
        y     = widget.value.layout[0].y
        w     = widget.value.layout[0].w
        h     = widget.value.layout[0].h
        min_h = widget.value.layout[0].min_h
      }
    }
  }
}
