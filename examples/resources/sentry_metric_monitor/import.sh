# import using the full URL:
terraform import sentry_metric_monitor.default https://{organization}.sentry.io/monitors/{id}/

# import using the organization and monitor id from the URL:
# https://{organization}.sentry.io/monitors/{id}/
terraform import sentry_metric_monitor.default {organization}/{id}