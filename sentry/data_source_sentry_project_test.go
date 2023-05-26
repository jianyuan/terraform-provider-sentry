package sentry

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSentryProjectDataSource_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	dn := "data.sentry_project.test"
	rn := "sentry_project.test"

    var projectID string

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryProjectConfig_basic(teamName, projectName),
				Check:  resource.ComposeTestCheckFunc(
                    testAccCheckSentryProjectExists(rn, &projectID),
                    resource.TestCheckResourceAttrPair(dn, "organization", rn, "organization"),
                    resource.TestCheckResourceAttrPair(dn, "slug", rn, "slug"),
                    resource.TestCheckResourceAttrPair(dn, "internal_id", rn, "internal_id"),
                    resource.TestCheckResourceAttrPair(dn, "is_public", rn, "is_public"),
                ),
			},
		},
	})
}

func testAccSentryProjectConfig_basic(teamName, projectName string) string {
	return testAccSentryProjectConfig_team(teamName, projectName) + `
data "sentry_project" "test" {
	organization = sentry_project.test.organization
	slug         = sentry_project.test.slug
}
	`
}
