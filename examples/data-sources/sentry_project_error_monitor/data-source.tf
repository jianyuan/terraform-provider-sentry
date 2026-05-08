data "sentry_project_error_monitor" "example" {
  organization = "my-org"     # Or Organization ID
  project      = "my-project" # Or Project ID
}

output "project_error_monitor_id" {
  value = data.sentry_project_error_monitor.example.id
}
