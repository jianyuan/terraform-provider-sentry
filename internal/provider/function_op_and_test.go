package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
)

func TestOpAndFunction_known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::op_and(
							provider::sentry::op_status_code_check("greater_than", 199),
							provider::sentry::op_status_code_check("less_than", 300),
							provider::sentry::op_and(
								provider::sentry::op_jsonpath(
									provider::sentry::op_jsonpath_operand_literal("value"),
									"equals",
									"$.status",
								),
							),
						)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"op":"and","children":[{"op":"status_code_check","operator":{"cmp":"greater_than"},"value":199},{"op":"status_code_check","operator":{"cmp":"less_than"},"value":300},{"op":"and","children":[{"op":"json_path","operand":{"jsonpath_op":"literal","value":"value"},"operator":{"cmp":"equals"},"value":"$.status"}]}]}`)),
					statecheck.ExpectKnownOutputValue("test", acctest.StringConformingJsonSchema(sentrydata.MustResolvedUptimeAssertionSchemaForDefinition("OpAnd"))),
				},
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_and(
							"bogus",
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "children" parameter: Invalid JSON String Value: A string value was provided that is not valid JSON string format (RFC 7159).`),
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_and(
							"{}",
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "children" parameter: validating root: validating /$defs/Op: oneOf: did not validate against any of [anchor OpAnd anchor OpOr anchor OpNot anchor OpStatusCode anchor OpHeaderCheck anchor OpJsonPath].`),
			},
		},
	})
}

func TestOpAnd_null(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::op_and(
							null,
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "children" parameter: argument must not be null.`),
			},
		},
	})
}

func TestOpAnd_unknown(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "terraform_data" "child_1" {
						input = provider::sentry::op_status_code_check("greater_than", 199)
					}

					resource "terraform_data" "child_2" {
						input = provider::sentry::op_status_code_check("less_than", 300)
					}

					output "test" {
						value = provider::sentry::op_and(
							terraform_data.child_1.output,
							terraform_data.child_2.output,
						)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"op":"and","children":[{"op":"header_check","key_op":{"cmp":"equals"},"key_operand":{"header_op":"literal","value":"X-Header-Key"},"value_op":{"cmp":"equals"},"value_operand":{"header_op":"literal","value":"header-value"}}]}`)),
					statecheck.ExpectKnownOutputValue("test", acctest.StringConformingJsonSchema(sentrydata.MustResolvedUptimeAssertionSchemaForDefinition("OpAnd"))),
				},
			},
		},
	})
}
