# Configure the Sentry Provider. Sentry will proxy most requests to the correct region based on the organization.
# To avoid the overhead of the proxy, you can configure the provider to use a specific region.
provider "sentry" {
  token = var.sentry_auth_token
}

# Configure the Sentry Provider for the US region
provider "sentry" {
  token = var.sentry_auth_token

  base_url = "https://us.sentry.io/api/"
}

# Configure the Sentry Provider for the EU region
provider "sentry" {
  token = var.sentry_auth_token

  base_url = "https://de.sentry.io/api/"
}

# Configure the Sentry Provider for self-hosted Sentry
provider "sentry" {
  token = var.sentry_auth_token

  # The URL format must be "https://[hostname]/api/".
  base_url = "https://example.com/api/"
}
