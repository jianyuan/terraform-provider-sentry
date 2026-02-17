resource "sentry_monitor" "main" {
  organization = sentry_project.main.organization
  project      = sentry_project.main.id
  name         = "nightly-db-backup"

  schedule_crontab = "0 2 * * *"
  timezone         = "UTC"
  checkin_margin = 5
  max_runtime    = 60

  failure_issue_threshold = 2
  recovery_threshold      = 1
}
