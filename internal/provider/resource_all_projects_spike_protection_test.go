package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccAllProjectsSpikeProtectionResource(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	project1Name := acctest.RandomWithPrefix("tf-project")
	project2Name := acctest.RandomWithPrefix("tf-project")
	rn := "sentry_all_projects_spike_protection.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectResourceConfig(teamName, project1Name) + `
					resource "sentry_all_projects_spike_protection" "test" {
						organization = sentry_team.test.organization
						projects     = [sentry_project.test.id]
						enabled      = true
					}
				`,
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("projects"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact(project1Name),
					})),
				),
			},
			{
				Config: testAccProjectResourceConfig(teamName, project1Name) + `
					resource "sentry_all_projects_spike_protection" "test" {
						organization = sentry_team.test.organization
						projects     = [sentry_project.test.id]
						enabled      = false
					}
				`,
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("projects"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact(project1Name),
					})),
				),
			},
			{
				Config: testAccProjectResourceConfig(teamName, project1Name) + fmt.Sprintf(`
					resource "sentry_project" "test2" {
						organization = sentry_team.test.organization
						teams        = [sentry_team.test.id]
						name         = "%[1]s"
					}

					resource "sentry_all_projects_spike_protection" "test" {
						organization = sentry_team.test.organization
						projects     = [sentry_project.test.id, sentry_project.test2.id]
						enabled      = false
					}
				`, project2Name),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("projects"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact(project1Name),
						knownvalue.StringExact(project2Name),
					})),
				),
			},
		},
	})
}
