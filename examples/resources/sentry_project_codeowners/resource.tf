resource "sentry_project_codeowners" "default" {
  organization    = "my-organization"
  project         = "my-project"
  code_mapping_id = "123"
  raw             = file("${path.module}/CODEOWNERS")
}
