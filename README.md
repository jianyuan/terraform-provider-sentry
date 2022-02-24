# terraform-provider-sentry

[![Go Report Card](https://goreportcard.com/badge/github.com/jianyuan/terraform-provider-sentry)](https://goreportcard.com/report/github.com/jianyuan/terraform-provider-sentry)

Terraform provider for [Sentry](https://sentry.io).

[This package is also published on the official Terraform registry](https://registry.terraform.io/providers/jianyuan/sentry/latest).

## Usage

[See the docs for usage information](./docs).

## Installation

See the [the Provider Configuration page of the Terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) for instructions.

Pre-compiled binaries are available from the [Releases](https://github.com/jianyuan/terraform-provider-sentry/releases) page.

## Development

If you wish to work on the provider, you will need to install [Go](https://go.dev/doc/install) (We use >= 1.17) on your machine.

### Test

In order to run the full suite of acceptance tests, run `make testacc`.

Make sure to set the following environment variables beforehand:

- `SENTRY_TEST_ORGANIZATION`
- `SENTRY_AUTH_TOKEN`

_Note:_ Acceptance tests create real resources, and often cost money to run.
