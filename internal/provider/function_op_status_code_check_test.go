package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
)

func TestOpStatusCodeCheckFunction_known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::op_status_code_check("greater_than", 199)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"op":"status_code_check","operator":{"cmp":"greater_than"},"value":199}`)),
					statecheck.ExpectKnownOutputValue("test", acctest.StringConformingJsonSchema(sentrydata.MustResolvedUptimeAssertionSchemaForDefinition("OpStatusCode"))),
				},
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_status_code_check("less_than", 300)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"op":"status_code_check","operator":{"cmp":"less_than"},"value":300}`)),
					statecheck.ExpectKnownOutputValue("test", acctest.StringConformingJsonSchema(sentrydata.MustResolvedUptimeAssertionSchemaForDefinition("OpStatusCode"))),
				},
			},
			{
				Config: `
					output "test" {
						value = provider::sentry::op_status_code_check("bogus", 199)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "operator" parameter: validating root: validating /$defs/ComparisonType: enum: bogus does not equal any of: [equals not_equal less_than greater_than always never].`),
			},
		},
	})
}

func TestOpStatusCodeCheckFunction_null(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::sentry::op_status_code_check(null, null)
					}
				`,
				ExpectError: acctest.ExpectLiteralError(`Invalid value for "operator" parameter: argument must not be null.`),
			},
		},
	})
}

func TestOpStatusCodeCheckFunction_unknown(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "terraform_data" "operator" {
						input = "greater_than"
					}

					resource "terraform_data" "value" {
						input = 199
					}

					output "test" {
						value = provider::sentry::op_status_code_check(terraform_data.operator.output, terraform_data.value.output)
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(`{"op":"status_code_check","operator":{"cmp":"greater_than"},"value":199}`)),
					statecheck.ExpectKnownOutputValue("test", acctest.StringConformingJsonSchema(sentrydata.MustResolvedUptimeAssertionSchemaForDefinition("OpStatusCode"))),
				},
			},
		},
	})
}
