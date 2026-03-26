package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
)

func TestOpHeaderOperandGlobFunction_known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_operand_glob("value*")
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"header_op":"glob","pattern":{"value":"value*"}}`)),
					statecheck.ExpectKnownOutputValue("test", acctest.StringConformingJsonSchema(sentrydata.MustResolvedUptimeAssertionSchemaForDefinition("OpHeaderOperandGlob"))),
				},
			},
		},
	})
}

func TestOpHeaderOperandGlobFunction_null(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::op_header_operand_glob(null)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "pattern" parameter: argument must not be null.`),
			},
		},
	})
}

func TestOpHeaderOperandGlobFunction_unknown(t *testing.T) {
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
						value = provider::sentry::op_header_operand_glob(terraform_data.value.output)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"header_op":"glob","pattern":{"value":"value"}}`)),
					statecheck.ExpectKnownOutputValue("test", acctest.StringConformingJsonSchema(sentrydata.MustResolvedUptimeAssertionSchemaForDefinition("OpHeaderOperandGlob"))),
				},
			},
		},
	})
}
