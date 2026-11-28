resource "sentry_organization_auth_token" "ci" {
  organization = "my-organization"
  name         = "GitHub Actions"
}

resource "aws_secretsmanager_secret" "sentry_auth_token" {
  name = "ci/sentry-auth-token"
}

resource "aws_secretsmanager_secret_version" "sentry_auth_token" {
  secret_id        = aws_secretsmanager_secret.sentry_auth_token.id
  secret_string_wo = sentry_organization_auth_token.ci.token
  # Pushes the new value whenever the token is replaced.
  secret_string_wo_version = tonumber(sentry_organization_auth_token.ci.id)
}
