# import using the organization slug and the mapping id from the Sentry API:
# GET /api/0/organizations/[org-slug]/external-users/
terraform import sentry_organization_user_mapping.jane_doe_gitlab org-slug/mapping-id
