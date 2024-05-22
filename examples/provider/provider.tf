# Configure the Sentry Provider for US data storage location (default)
provider "sentry" {
  token = var.sentry_auth_token

  # If you want to be explicit, you can specify the base URL for the US region.
  # base_url = "https://us.sentry.io/api/"
  # or
  # base_url = "https://sentry.io/api/"
}

# Configure the Sentry Provider for EU data storage location
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
