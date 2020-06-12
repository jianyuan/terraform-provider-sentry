# terraform-provider-sentry

[![CircleCI](https://circleci.com/gh/jianyuan/terraform-provider-sentry/tree/master.svg?style=svg)](https://circleci.com/gh/jianyuan/terraform-provider-sentry/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/jianyuan/terraform-provider-sentry)](https://goreportcard.com/report/github.com/jianyuan/terraform-provider-sentry)

Terraform provider for [Sentry](https://sentry.io).

[This package is also published on the official Terraform registry](https://registry.terraform.io/providers/jianyuan/sentry/latest).

## Usage

[See the docs for usage information](./docs).

## Installation

See the [the Provider Configuration page of the Terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) for instructions.

Pre-compiled binaries are available from the [Releases](https://github.com/jianyuan/terraform-provider-sentry/releases) page.

## Development

### Test

Test the provider by running `make test`.

Make sure to set the following environment variables:

- `SENTRY_TEST_ORGANIZATION`
- `SENTRY_TOKEN`

### Build

See the [Writing Custom Providers page of the Terraform documentation](https://www.terraform.io/docs/extend/writing-custom-providers.html#building-the-plugin) for instructions.
