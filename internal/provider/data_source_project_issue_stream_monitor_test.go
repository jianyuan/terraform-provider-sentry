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

func TestAccProjectIssueStreamMonitorDataSource_basic(t *testing.T) {
	rn := "data.sentry_project_issue_stream_monitor.test"

	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("organization"), knownvalue.StringExact(acctest.TestOrganization)),
		statecheck.ExpectKnownValue(rn, tfjsonpath.New("project"), knownvalue.NotNull()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectIssueStreamMonitorDataSourceConfig(),
				ConfigStateChecks: append(
					checks,
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.StringExact("Issue Stream")),
					statecheck.ExpectKnownValue(rn, tfjsonpath.New("owner"), knownvalue.Null()),
				),
			},
		},
	})
}

func testAccProjectIssueStreamMonitorDataSourceConfig() string {
	return fmt.Sprintf(`
		data "sentry_project_issue_stream_monitor" "test" {
			organization = "%s"
			project      = "%s"
		}
	`, acctest.TestOrganization, acctest.TestProject.Slug)
}
