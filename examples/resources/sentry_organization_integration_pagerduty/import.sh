# import using the organization slug from the URL:
# https://sentry.io/api/0/organizations/[org-slug]/integrations/
# [integration-id] is the top-level `id` of the PagerDuty organization integration
# [service-id] is the `id` of the service_table record to import under the configData property
terraform import sentry_organization_integration_pagerduty.this org-slug/integration-id/service-id