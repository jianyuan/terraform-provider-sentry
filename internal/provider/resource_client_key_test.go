package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func TestAccClientKeyResource_MigrateFromPluginSDK(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")
	rn := "sentry_key.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_csp"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					acctest.ProviderName: {
						Source:            "jianyuan/sentry",
						VersionConstraint: "0.12.3",
					},
				},
				Config:            testAccClientKeyResourceConfig(teamName, projectName, keyName, ""),
				ConfigStateChecks: checks,
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   testAccClientKeyResourceConfig(teamName, projectName, keyName, ""),
				ConfigStateChecks:        checks,
			},
		},
	})
}

func TestAccClientKeyResource(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")
	rn := "sentry_key.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_csp"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClientKeyResourceConfig(teamName, projectName, keyName, ""),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
				),
			},
			{
				Config: testAccClientKeyResourceConfig(teamName, projectName, keyName+"-renamed", ""),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
				),
			},
			{
				Config: testAccClientKeyResourceConfig(teamName, projectName, keyName+"-renamed", `
					rate_limit_window = 1
					rate_limit_count  = 2
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Int64Exact(2)),
				),
			},
			{
				Config: testAccClientKeyResourceConfig(teamName, projectName, keyName+"-renamed", `
					rate_limit_window = 3
					rate_limit_count  = 4
				`),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Int64Exact(4)),
				),
			},
			{
				Config: testAccClientKeyResourceConfig(teamName, projectName, keyName+"-renamed", ""),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName+"-renamed")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
				),
			},
			{
				ResourceName: rn,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[rn]
					if !ok {
						return "", fmt.Errorf("not found: %s", rn)
					}
					organization := rs.Primary.Attributes["organization"]
					project := rs.Primary.Attributes["project"]
					keyId := rs.Primary.ID
					return buildThreePartID(organization, project, keyId), nil
				},
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClientKeyResourceConfig(teamName, projectName, keyName, extras string) string {
	return testAccProjectResourceConfig(teamName, projectName) + fmt.Sprintf(`
resource "sentry_key" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	name         = "%[1]s"
	%[2]s
}
`, keyName, extras)
}
