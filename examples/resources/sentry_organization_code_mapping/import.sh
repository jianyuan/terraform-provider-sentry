# import using the organization slug from the URL:
# https://sentry.io/settings/[org-slug]/integrations/github/[org-integration-id]/
# and inspect network tab for request to https://sentry.io/api/0/organizations/[org-slug]/code-mappings/
# find the corresponding list element and reference [code-mapping-id] from the key "id"
terraform import sentry_organization_code_mapping.this org-slug/31347
