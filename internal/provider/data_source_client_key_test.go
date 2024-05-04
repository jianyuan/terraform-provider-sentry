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

func TestAccClientKeyDataSource_MigrateFromPluginSDK(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	rn := "data.sentry_key.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("Default")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("first"), knownvalue.Bool(false)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
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
				Config:            testAccClientKeyDataSourceConfig_bare(teamName, projectName),
				ConfigStateChecks: checks,
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   testAccClientKeyDataSourceConfig_bare(teamName, projectName),
				ConfigStateChecks:        checks,
			},
		},
	})
}

func TestAccClientKeyDataSource_id(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")
	rn := "data.sentry_key.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("first"), knownvalue.Bool(false)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_csp"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClientKeyDataSourceConfig(teamName, projectName, keyName, `
					id = sentry_key.test.id
				`),
				ConfigStateChecks: checks,
			},
		},
	})
}

func TestAccClientKeyDataSource_name(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")
	rn := "data.sentry_key.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact(keyName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("first"), knownvalue.Bool(false)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_csp"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClientKeyDataSourceConfig(teamName, projectName, keyName, `
					name = sentry_key.test.name
				`),
				ConfigStateChecks: checks,
			},
		},
	})
}

func TestAccClientKeyDataSource_first(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	keyName := acctest.RandomWithPrefix("tf-key")
	rn := "data.sentry_key.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.StringExact(projectName)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("Default")),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("first"), knownvalue.Bool(true)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project_id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_window"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("rate_limit_count"), knownvalue.Null()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_public"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_secret"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("dsn_csp"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClientKeyDataSourceConfig(teamName, projectName, keyName, `
					first = true
				`),
				ConfigStateChecks: checks,
			},
		},
	})
}

func testAccClientKeyDataSourceConfig_bare(teamName, projectName string) string {
	return testAccProjectResourceConfig(teamName, projectName) + `
data "sentry_key" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
}
`
}

func testAccClientKeyDataSourceConfig(teamName, projectName, keyName, extras string) string {
	return testAccClientKeyResourceConfig(teamName, projectName, keyName, "") + fmt.Sprintf(`
data "sentry_key" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	%[1]s
}
`, extras)
}
