data "sentry_project_issue_stream_monitor" "example" {
  organization = "my-org"     # Or Organization ID
  project      = "my-project" # Or Project ID
}

output "project_issue_stream_monitor_id" {
  value = data.sentry_project_issue_stream_monitor.example.id
}
