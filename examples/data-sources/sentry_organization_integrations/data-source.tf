# Example 1: Simple Map from Bulk Lookup
# This is the most straightforward method, using the locals block to create a map.
# This creates a map where keys are integration names and values are integration IDs.
# This approach assumes the combination of provider_key and name is unique for each integration.
data "sentry_organization_integrations" "all_github" {
  organization = "my-organization"
  provider_key = "github"
}

locals {
  # Map structure: { "org-name-1" = "id-1", "org-name-2" = "id-2" }
  github_integration_map = {
    for integration in data.sentry_organization_integrations.all_github.integrations :
    integration.name => integration.id
  }
}

# Example usage elsewhere in code:
# local.github_integration_map["my-github-org"]


# Example 2: Map across All Providers (Omni-lookup)
# This approach is useful when you need to look up integrations across multiple providers.
data "sentry_organization_integrations" "all" {
  organization = "my-organization"
  # provider_key omitted - returns all integrations for the organization
}

locals {
  # Map structure: { "github:my-org" = "id-1", "slack:my-workspace" = "id-2" }
  all_integrations_map = {
    for i in data.sentry_organization_integrations.all.integrations :
    "${i.provider_key}:${i.name}" => i.id
  }
}

# Example usage elsewhere in code:
# local.all_integrations_map["github:my-github-org"]


# Example 3: Specific Integration Lookup
# This approach is useful when you need to look up a specific integration by name.
# The `one()` function ensures that only one integration is returned.
data "sentry_organization_integrations" "all_github" {
  organization = "my-organization"
  provider_key = "github"
}

locals {
  # Match a specific integration from the bulk list
  matched_integration = one([
    for i in data.sentry_organization_integrations.all_github.integrations :
    i if i.name == "my-github-organization"
  ])
}

output "github_id" {
  # Safely returns the ID if found, or null if no match exists
  value = try(local.matched_integration.id, null)
}
