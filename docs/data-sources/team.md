---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sentry_team Data Source - terraform-provider-sentry"
subcategory: ""
description: |-
  Sentry Team data source.
---

# sentry_team (Data Source)

Sentry Team data source.

## Example Usage

```terraform
# Retrieve a team
data "sentry_team" "default" {
  organization = "my-organization"

  slug = "my-team"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `organization` (String) The slug of the organization the team should be created for.
- `slug` (String) The unique URL slug for this team.

### Read-Only

- `has_access` (Boolean)
- `id` (String) The ID of this resource.
- `internal_id` (String) The internal ID for this team.
- `is_member` (Boolean)
- `is_pending` (Boolean)
- `name` (String) The human readable name for this organization.
