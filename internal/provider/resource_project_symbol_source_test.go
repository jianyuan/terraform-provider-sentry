package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/resourceid"
)

func TestAccProjectSymbolSourceResource(t *testing.T) {
	rn := "sentry_project_symbol_source.test"
	project := acctest.RandomWithPrefix("tf-project")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectSymbolSourceConfig(project, "s3"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", project),
					resource.TestCheckResourceAttr(rn, "type", "s3"),
					resource.TestCheckResourceAttr(rn, "name", "s3"),
					resource.TestCheckResourceAttr(rn, "layout.%", "2"),
					resource.TestCheckResourceAttr(rn, "layout.type", "native"),
					resource.TestCheckResourceAttr(rn, "layout.casing", "default"),
					resource.TestCheckResourceAttr(rn, "bucket", "bucket"),
					resource.TestCheckResourceAttr(rn, "region", "us-east-1"),
					resource.TestCheckResourceAttr(rn, "access_key", "access_key"),
					resource.TestCheckResourceAttr(rn, "secret_key", "secret_key"),
				),
			},
			{
				Config: testAccProjectSymbolSourceConfig(project, "s3-edited"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
					resource.TestCheckResourceAttr(rn, "project", project),
					resource.TestCheckResourceAttr(rn, "type", "s3"),
					resource.TestCheckResourceAttr(rn, "name", "s3-edited"),
					resource.TestCheckResourceAttr(rn, "layout.%", "2"),
					resource.TestCheckResourceAttr(rn, "layout.type", "native"),
					resource.TestCheckResourceAttr(rn, "layout.casing", "default"),
					resource.TestCheckResourceAttr(rn, "bucket", "bucket"),
					resource.TestCheckResourceAttr(rn, "region", "us-east-1"),
					resource.TestCheckResourceAttr(rn, "access_key", "access_key"),
					resource.TestCheckResourceAttr(rn, "secret_key", "secret_key"),
				),
			},
			{
				ResourceName:            rn,
				ImportState:             true,
				ImportStateIdFunc:       resourceid.ImportState3PartIDFunc(rn, "organization", "project", "id"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
		},
	})
}

func testAccProjectSymbolSourceConfig(projectName string, name string) string {
	return fmt.Sprintf(`
resource "sentry_project" "test" {
	organization = "%[1]s"
	teams        = ["%[2]s"]
	name         = "%[3]s"
	platform     = "go"
}

resource "sentry_project_symbol_source" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	type         = "s3"
	name         = "%[4]s"
	layout       = {
		type   = "native"
		casing = "default"
	}
	bucket       = "bucket"
	region       = "us-east-1"
	access_key   = "access_key"
	secret_key   = "secret_key"
}
`, acctest.TestOrganization, acctest.TestTeam.Slug, projectName, name)
}

func TestAccProjectSymbolSourceResource_writeOnly(t *testing.T) {
	rn := "sentry_project_symbol_source.test"
	project := acctest.RandomWithPrefix("tf-project")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { acctest.PreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_11_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectSymbolSourceWriteOnlyConfig(project, "s3", "secret_key_1", 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("s3")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key_wo"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key_wo_version"), knownvalue.Int64Exact(1)),
				},
			},
			// Updating another attribute keeps the stored secret.
			{
				Config: testAccProjectSymbolSourceWriteOnlyConfig(project, "s3-edited", "secret_key_1", 1),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("s3-edited")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key"), knownvalue.Null()),
				},
			},
			// Bumping the version sends the new secret.
			{
				Config: testAccProjectSymbolSourceWriteOnlyConfig(project, "s3-edited", "secret_key_2", 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key"), knownvalue.Null()),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("secret_key_wo_version"), knownvalue.Int64Exact(2)),
				},
			},
			{
				ResourceName:            rn,
				ImportState:             true,
				ImportStateIdFunc:       resourceid.ImportState3PartIDFunc(rn, "organization", "project", "id"),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key_wo", "secret_key_wo_version"},
			},
		},
	})
}

func testAccProjectSymbolSourceWriteOnlyConfig(projectName, name, secretKey string, secretKeyVersion int) string {
	return fmt.Sprintf(`
resource "sentry_project" "test" {
	organization = "%[1]s"
	teams        = ["%[2]s"]
	name         = "%[3]s"
	platform     = "go"
}

resource "sentry_project_symbol_source" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	type         = "s3"
	name         = "%[4]s"
	layout       = {
		type   = "native"
		casing = "default"
	}
	bucket                = "bucket"
	region                = "us-east-1"
	access_key            = "access_key"
	secret_key_wo         = "%[5]s"
	secret_key_wo_version = %[6]d
}
`, acctest.TestOrganization, acctest.TestTeam.Slug, projectName, name, secretKey, secretKeyVersion)
}
