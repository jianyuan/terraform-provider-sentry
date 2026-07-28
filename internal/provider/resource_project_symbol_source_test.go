package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/resourceid"
)

func TestAccProjectSymbolSourceResource(t *testing.T) {
	rn := "sentry_project_symbol_source.test"
	project := acctest.RandomWithPrefix("tf-project")

	resource.Test(t, resource.TestCase{
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
