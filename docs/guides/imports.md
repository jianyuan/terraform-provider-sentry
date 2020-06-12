# Import existing resources

You can import existing resources using the [`terraform import`](https://www.terraform.io/docs/import/usage.html) command.

To import an organization:

```bash
$ terraform import sentry_organization.default org-slug
```

To import a team:

```bash
$ terraform import sentry_team.default org-slug/team-slug
```

To import a project:

```bash
$ terraform import sentry_project.default org-slug/project-slug
```
