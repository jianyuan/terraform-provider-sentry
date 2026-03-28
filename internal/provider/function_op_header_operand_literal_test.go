package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
)

func TestOpHeaderOperandLiteralFunction_known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_operand_literal("value")
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"header_op":"literal","value":"value"}`)),
					statecheck.ExpectKnownOutputValue("test", acctest.StringConformingJsonSchema(sentrydata.MustResolvedUptimeAssertionSchemaForDefinition("OpHeaderOperandLiteral"))),
				},
			},
		},
	})
}

func TestOpHeaderOperandLiteralFunction_null(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_operand_literal(null)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "value" parameter: argument must not be null.`),
			},
		},
	})
}

func TestOpHeaderOperandLiteralFunction_unknown(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "terraform_data" "value" {
						input = "value"					
					}

					output "test" {
						value = provider::sentry::op_header_operand_literal(terraform_data.value.output)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"header_op":"literal","value":"value"}`)),
					statecheck.ExpectKnownOutputValue("test", acctest.StringConformingJsonSchema(sentrydata.MustResolvedUptimeAssertionSchemaForDefinition("OpHeaderOperandLiteral"))),
				},
			},
		},
	})
}
