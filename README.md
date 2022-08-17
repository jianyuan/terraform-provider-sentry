# Terraform Provider Sentry

[![Go Report Card](https://goreportcard.com/badge/github.com/jianyuan/terraform-provider-sentry)](https://goreportcard.com/report/github.com/jianyuan/terraform-provider-sentry)

<a href="https://sentry.io/?utm_source=terraform&utm_medium=docs" target="_blank">
    <img src="templates/sentry-wordmark-dark-280x84.svg" alt="Sentry" width="280">
</a>

<a href="https://www.terraform.io/" target="_blank">
    <img src="templates/Terraform_PrimaryLogo_Color_RGB.svg" alt="Terraform" width="280">
</a>

The Terraform provider for [Sentry](https://sentry.io/?utm_source=terraform&utm_medium=docs) allows teams to configure and update Sentry project parameters via their command line. This provider is officially sponsored by [Sentry](https://sentry.io/?utm_source=terraform&utm_medium=docs).

## Usage

Detailed documentation is available on the [Terraform provider registry](https://registry.terraform.io/providers/jianyuan/sentry/latest).

## Development

If you wish to work on the provider, you will need to install [Go](https://go.dev/doc/install) (We use >= 1.18) on your machine.

### Test

In order to run the full suite of acceptance tests, run `make testacc`.

Make sure to set the following environment variables beforehand:

- `SENTRY_TEST_ORGANIZATION`
- `SENTRY_AUTH_TOKEN`

_Note:_ Acceptance tests create real resources, and often cost money to run.
