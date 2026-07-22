# import using the full URL:
terraform import sentry_alert.default https://{organization}.sentry.io/monitors/alerts/{id}/

# import using the organization and alert id from the URL:
# https://{organization}.sentry.io/monitors/alerts/{id}/
terraform import sentry_alert.default organization/id
