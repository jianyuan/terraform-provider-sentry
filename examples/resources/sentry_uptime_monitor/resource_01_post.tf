# POST request with body and headers
resource "sentry_uptime_monitor" "test" {
  organization = data.sentry_organization.test.slug
  project      = sentry_project.test.slug
  name         = "Uptime check for sentry.io"

  owner = {
    team_id = sentry_team.test.internal_id
  }

  environment = "production"

  url              = "https://sentry.io"
  method           = "POST"
  body             = <<EOT
    {
      "key": "value"
    }
  EOT
  interval_seconds = 60
  timeout_ms       = 5000
  headers = {
    "X-Header-Key" : "X-Header-Value"
  }

  # Assertion that checks for a 2xx status code response.
  assertion_json = provider::sentry::assertion(
    provider::sentry::op_and(
      provider::sentry::op_status_code_check("greater_than", 199),
      provider::sentry::op_status_code_check("less_than", 300),
    )
  )
  # The following is equivalent to the above.
  # assertion_json = <<EOT
  #   {
  #     "root": {
  #       "op": "and",
  #       "children": [
  #         {"op": "status_code_check", "operator": {"cmp": "greater_than"}, "value": 199},
  #         {"op": "status_code_check", "operator": {"cmp": "less_than"}, "value": 300}
  #       ]
  #     }
  #   }
  # EOT
}
