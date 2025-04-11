
# sentry_cron_monitor

Manages a Sentry Cron Monitor.

## Example Usage

```hcl
resource "sentry_cron_monitor" "example" {
  organization        = "my-org"
  project             = "my-project"
  name                = "daily-cron"
  schedule_type       = "crontab"
  schedule            = "0 0 * * *"
  check_in_margin     = 5
  max_runtime         = 30
  timezone            = "UTC"
  alert_rule_enabled  = true
  alert_threshold     = 1
}
```
