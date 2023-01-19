# import using the organization, project slugs and rule id from the URL:
# https://sentry.io/organizations/[org-slug]/projects/[project-slug]/
# https://sentry.io/organizations/[org-slug]/alerts/rules/details/[rule-id]/
# or
# https://sentry.io/organizations/[org-slug]/alerts/metric-rules/[project-slug]/[rule-id]/
terraform import sentry_metric_alert.default org-slug/project-slug/rule-id
