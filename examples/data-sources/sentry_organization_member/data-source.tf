data "sentry_organization_member" "default" {
  organization = "terraform-provider-sentry"
  email        = "test@example.com"
}
