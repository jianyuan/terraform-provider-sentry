# Retrieve a project key by id
data "sentry_key" "default" {
  organization = "my-organization"
  project      = "web-app"

  id = "73e6e1c04501397c0f87f36bf48f22ea"
}

# Retrieve a project key by name
data "sentry_key" "default" {
  organization = "my-organization"
  project      = "web-app"

  name = "Default"
}

# Retrieve the first key of a project
data "sentry_key" "first" {
  organization = "my-organization"
  project      = "web-app"

  first = true
}
