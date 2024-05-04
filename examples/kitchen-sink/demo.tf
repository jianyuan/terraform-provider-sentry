terraform {
  required_providers {
    sentry = {
      source = "jianyuan/sentry"
    }
  }
}

provider "sentry" {
  token = "xxx"
}

data "sentry_organization" "main" {
  # Taken from URL: https://sentry.io/organizations/[slug]/issues/
  slug = "my-sentry-organization"
}

output "organization" {
  value = data.sentry_organization.main
}

#
# Team
#

resource "sentry_team" "main" {
  organization = data.sentry_organization.main.id
  name         = "My team"
}

output "team" {
  value = sentry_team.main.id
}

#
# Project
#

resource "sentry_project" "main" {
  organization = sentry_team.main.organization
  team         = sentry_team.main.id
  name         = "My project"
  platform     = "python"
}

output "project" {
  value = sentry_project.main.id
}

#
# Client key
#

data "sentry_key" "main" {
  organization = sentry_project.main.organization
  project      = sentry_project.main.id

  first = true
}

output "client_key" {
  value = data.sentry_key.main
}
