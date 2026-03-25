package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestOpJsonpathFunction_known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::op_jsonpath(
							provider::sentry::op_jsonpath_operand_literal("value"),
							"equals",
							"$.status",
						)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"op":"json_path","operand":{"jsonpath_op":"literal","value":"value"},"operator":{"cmp":"equals"},"value":"$.status"}`)),
				},
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_jsonpath(
							provider::sentry::op_jsonpath_operand_glob("value*"),
							"equals",
							"$.status",
						)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"op":"json_path","operand":{"jsonpath_op":"glob","pattern":{"value":"value*"}},"operator":{"cmp":"equals"},"value":"$.status"}`)),
				},
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_jsonpath(
							"bogus",
							"equals",
							"$.status",
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "operand" parameter: Invalid JSON String Value: A string value was provided that is not valid JSON string format (RFC 7159).`),
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_jsonpath(
							"{}",
							"equals",
							"$.status",
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "operand" parameter: oneOf: Value does not match the oneOf schema.`),
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_jsonpath(
							provider::sentry::op_jsonpath_operand_literal("value"),
							"bogus",
							"$.status",
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "operator" parameter: enum: Value bogus should be one of the allowed values: equals, not_equal, less_than, greater_than, always, never.`),
			},
		},
	})
}

func TestOpJsonpathFunction_null(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::op_jsonpath(null, null, null)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "operand" parameter: argument must not be null.`),
			},
		},
	})
}

func TestOpJsonpathFunction_unknown(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "terraform_data" "operand" {
						input = provider::sentry::op_jsonpath_operand_literal("value")
					}

					resource "terraform_data" "operator" {
						input = "equals"
					}

					resource "terraform_data" "value" {
						input = "$.status"
					}

					output "test" {
						value = provider::sentry::op_jsonpath(
							terraform_data.operand.output,
							terraform_data.operator.output,
							terraform_data.value.output,
						)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"op":"json_path","operand":{"jsonpath_op":"literal","value":"value"},"operator":{"cmp":"equals"},"value":"$.status"}`)),
				},
			},
		},
	})
}
