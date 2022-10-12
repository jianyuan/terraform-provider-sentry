# sentry_organization_member Resource

Sentry Organization Member resource.

## Example Usage

```hcl
# Create an organization member
resource "sentry_organization_member" "john_doe" {
  email = "test@example.com"
  role  = "member"
  teams = ["my-team"]
}
```

## Argument Reference

The following arguments are supported:

- `email` - (Required) The email address of the user you want to invite.
- `role` - (Required) The role of the user you want to invite. Must be one of the following: `member`, `manager`, `billing`, `admin` or `owner`. 
- `teams` - (Optional) The list of slugs for each team you want to add the member to.

## Attribute Reference

The following attributes are exported:

- `member_id` - The ID of the created organization member.
- `teams` - The list of teams that the member is a part of.
- `role` - The role that the user has.
- `email` - The email address of the organization member.
- `pending` - The boolean acceptance status of the membership invite. e.g. True means the invite has not been accepted.
- `expired` - A boolean that tells you whether the membership invite has expired.

## Import

This resource can be imported using an ID made up of the organization slug and member ID.

```bash
$ terraform import sentry_organization.default org-slug/member-id
```
