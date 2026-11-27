# Retrieve a project key by name without storing it in state
ephemeral "sentry_key" "default" {
  organization = "my-organization"
  project      = "web-app"

  name = "Default"
}

# Store the DSN in AWS Secrets Manager using a write-only attribute
resource "aws_secretsmanager_secret" "sentry_dsn" {
  name = "web-app/sentry-dsn"
}

resource "aws_secretsmanager_secret_version" "sentry_dsn" {
  secret_id                = aws_secretsmanager_secret.sentry_dsn.id
  secret_string_wo         = ephemeral.sentry_key.default.dsn["public"]
  secret_string_wo_version = 1
}
