# Issue Stream Monitor: The default monitor tracking new issues of all types created for a project
data "sentry_monitor" "issue_stream" {
  organization = "my-org"     # Or Organization ID
  project      = "my-project" # Or Project ID

  type = "issue_stream"
}

output "issue_stream_monitor_id" {
  value = data.sentry_monitor.issue_stream.id
}
