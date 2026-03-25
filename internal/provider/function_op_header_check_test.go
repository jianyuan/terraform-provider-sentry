package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestOpHeaderCheckFunction_known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_check(
							"equals",
							provider::sentry::op_header_operand_literal("X-Header-Key"),
							"equals",
							provider::sentry::op_header_operand_literal("header-value"),
						)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"op":"header_check","key_op":{"cmp":"equals"},"key_operand":{"header_op":"literal","value":"X-Header-Key"},"value_op":{"cmp":"equals"},"value_operand":{"header_op":"literal","value":"header-value"}}`)),
				},
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_check(
							"bogus",
							provider::sentry::op_header_operand_literal("X-Header-Key"),
							"equals",
							provider::sentry::op_header_operand_literal("header-value"),
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "key_operator" parameter: enum: Value bogus should be one of the allowed values: equals, not_equal, less_than, greater_than, always, never.`),
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_check(
							"equals",
							provider::sentry::op_header_operand_literal("X-Header-Key"),
							"bogus",
							provider::sentry::op_header_operand_literal("header-value"),
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "value_operator" parameter: enum: Value bogus should be one of the allowed values: equals, not_equal, less_than, greater_than, always, never.`),
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_check(
							"equals",
							"bogus",
							"equals",
							provider::sentry::op_header_operand_literal("header-value"),
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "key_operand" parameter: Invalid JSON String Value: A string value was provided that is not valid JSON string format (RFC 7159).`),
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_check(
							"equals",
							provider::sentry::op_header_operand_literal("X-Header-Key"),
							"equals",
							"bogus",
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "value_operand" parameter: Invalid JSON String Value: A string value was provided that is not valid JSON string format (RFC 7159).`),
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_check(
							"equals",
							"{}",
							"equals",
							provider::sentry::op_header_operand_glob("header-value"),
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "key_operand" parameter: oneOf: Value does not match the oneOf schema.`),
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_check(
							"equals",
							provider::sentry::op_header_operand_literal("X-Header-Key"),
							"equals",
							"{}",
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "value_operand" parameter: oneOf: Value does not match the oneOf schema.`),
			},
		},
	})
}

func TestOpHeaderCheckFunction_null(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_check(
							null,
							null,
							null,
							null,
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "key_operator" parameter: argument must not be null.`),
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_check(
							"equals",
							null,
							null,
							null,
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "key_operand" parameter: argument must not be null.`),
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_check(
							"equals",
							provider::sentry::op_header_operand_literal("X-Header-Key"),
							null,
							null,
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "value_operator" parameter: argument must not be null.`),
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_check(
							"equals",
							provider::sentry::op_header_operand_literal("X-Header-Key"),
							"equals",
							null,
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "value_operand" parameter: argument must not be null.`),
			},
		},
	})
}

func TestOpHeaderCheckFunction_unknown(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "terraform_data" "key_operator" {
						input = "equals"
					}

					resource "terraform_data" "key_operand" {
						input = provider::sentry::op_header_operand_literal("X-Header-Key")
					}

					resource "terraform_data" "value_operator" {
						input = "equals"
					}

					resource "terraform_data" "value_operand" {
						input = provider::sentry::op_header_operand_literal("header-value")
					}

					output "test" {
						value = provider::sentry::op_header_check(
							terraform_data.key_operator.output,
							terraform_data.key_operand.output,
							terraform_data.value_operator.output,
							terraform_data.value_operand.output,
						)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"op":"header_check","key_op":{"cmp":"equals"},"key_operand":{"header_op":"literal","value":"X-Header-Key"},"value_op":{"cmp":"equals"},"value_operand":{"header_op":"literal","value":"header-value"}}`)),
				},
			},
		},
	})
}
