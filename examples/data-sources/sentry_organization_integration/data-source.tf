# Retrieve a Github organization integration
data "sentry_organization_integration" "github" {
  organization = "my-organization"

  provider_key = "github"
  name         = "my-github-organization" # Name of your GitHub organization (i.e. http://github.com/[name])
}

# Retrieve a Slack integration
data "sentry_organization_integration" "slack" {
  organization = "my-organization"

  provider_key = "slack"
  name         = "Slack Workspace" # Name of your Slack workspace
}
