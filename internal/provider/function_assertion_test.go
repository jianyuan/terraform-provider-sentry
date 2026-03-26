package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
)

func TestAssertionFunction_known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::assertion(
							provider::sentry::op_and(
								provider::sentry::op_status_code_check("greater_than", 199),
								provider::sentry::op_status_code_check("less_than", 300),
							)
						)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"root":{"op":"and","children":[{"op":"status_code_check","operator":{"cmp":"greater_than"},"value":199},{"op":"status_code_check","operator":{"cmp":"less_than"},"value":300}]}}`)),
					statecheck.ExpectKnownOutputValue("test", acctest.StringConformingJsonSchema(sentrydata.MustResolvedUptimeAssertionSchemaForDefinition("Assertion"))),
				},
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::assertion(
							"bogus",
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "root" parameter: Invalid JSON String Value: A string value was provided that is not valid JSON string format (RFC 7159).`),
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::assertion(
							"{}",
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "root" parameter: validating root: validating /$defs/Op: oneOf: did not validate against any of [anchor OpAnd anchor OpOr anchor OpNot anchor OpStatusCode anchor OpHeaderCheck anchor OpJsonPath].`),
			},
		},
	})
}

func TestAssertion_null(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::assertion(
							null,
						)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "root" parameter: argument must not be null.`),
			},
		},
	})
}

func TestAssertion_unknown(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "terraform_data" "root" {
						input = provider::sentry::op_and(
							provider::sentry::op_status_code_check("greater_than", 199),
							provider::sentry::op_status_code_check("less_than", 300),
						)
					}

					output "test" {
						value = provider::sentry::assertion(
							terraform_data.root.output,
						)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"root":{"op":"and","children":[{"op":"status_code_check","operator":{"cmp":"greater_than"},"value":199},{"op":"status_code_check","operator":{"cmp":"less_than"},"value":300}]}}`)),
					statecheck.ExpectKnownOutputValue("test", acctest.StringConformingJsonSchema(sentrydata.MustResolvedUptimeAssertionSchemaForDefinition("Assertion"))),
				},
			},
		},
	})
}
