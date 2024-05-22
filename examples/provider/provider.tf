# Configure the Sentry Provider
provider "sentry" {
  token = var.sentry_auth_token

  # If you are self-hosting Sentry, set the base URL here.
  # The URL format must be "https://[hostname]/api/".
  # base_url = "https://example.com/api/"
}
