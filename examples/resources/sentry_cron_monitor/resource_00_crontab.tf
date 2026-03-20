# Cron monitor with crontab schedule
resource "sentry_cron_monitor" "crontab" {
  organization = data.sentry_organization.default.slug
  project      = sentry_project.default.slug

  name = "My cron monitor"

  default_assignee = {
    team_id = sentry_team.default.internal_id
  }

  checkin_margin          = 1
  failure_issue_threshold = 1
  max_runtime             = 30
  recovery_threshold      = 30

  schedule = {
    crontab = "0 0 * * *"
  }
}
