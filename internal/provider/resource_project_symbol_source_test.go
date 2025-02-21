package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

func TestAccProjectSymbolSourceResource(t *testing.T) {
	rn := "sentry_project_symbol_source.test"
	team := acctest.RandomWithPrefix("tf-team")
	project := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectSymbolSourceConfig(team, project, "s3"),
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
				Config: testAccProjectSymbolSourceConfig(team, project, "s3-edited"),
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
				ResourceName: rn,
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[rn]
					if !ok {
						return "", fmt.Errorf("not found: %s", rn)
					}
					organization := rs.Primary.Attributes["organization"]
					project := rs.Primary.Attributes["project"]
					sourceId := rs.Primary.ID
					return tfutils.BuildThreePartId(organization, project, sourceId), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_key"},
			},
		},
	})
}

func testAccProjectSymbolSourceConfig(teamName string, projectName string, name string) string {
	return testAccOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_team" "test" {
	organization = data.sentry_organization.test.id
	name         = "%[1]s"
	slug         = "%[1]s"
}

resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.id]
	name         = "%[2]s"
	platform     = "go"
}

resource "sentry_project_symbol_source" "test" {
	organization = sentry_project.test.organization
	project      = sentry_project.test.id
	type         = "s3"
	name         = "%[3]s"
	layout       = {
		type   = "native"
		casing = "default"
	}
	bucket       = "bucket"
	region       = "us-east-1"
	access_key   = "access_key"
	secret_key   = "secret_key"
}
`, teamName, projectName, name)
}
