package sentry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func TestAccSentryDashboard_basic(t *testing.T) {
	var dashboard sentry.Dashboard

	dashboardTitle := acctest.RandomWithPrefix("tf-dashboard")
	rn := "sentry_dashboard.test"

	check := func(alertName string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckSentryDashboardExists(rn, &dashboard),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSentryIssueAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryDashboardConfig(dashboardTitle),
				Check:  check(dashboardTitle),
			},
			{
				Config: testAccSentryDashboardConfig(dashboardTitle + "-renamed"),
				Check:  check(dashboardTitle + "-renamed"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSentryDashboardExists(n string, dashboard *sentry.Dashboard) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		org, id, err := splitSentryDashboardID(rs.Primary.ID)
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*sentry.Client)
		ctx := context.Background()
		gotDashboard, _, err := client.Dashboards.Get(ctx, org, id)
		if err != nil {
			return err
		}
		*dashboard = *gotDashboard
		return nil
	}
}

func testAccSentryDashboardConfig(dashboardTitle string) string {
	return fmt.Sprintf(`
data "sentry_organization" "test" {
	slug = "%[1]s"
}

resource "sentry_dashboard" "test" {
	organization = data.sentry_organization.test.id
	title        = "%[2]s"

	widget {
		title        = "Custom Widget"
		display_type = "world_map"

		query {
			name       = "Metric"

			fields     = ["count()"]
			aggregates = ["count()"]
			conditions = "!event.type:transaction"
		}

		layout {
			x     = 0
			y     = 0
			w     = 2
			h     = 1
			min_h = 1
		}
	}
}
	`, testOrganization, dashboardTitle)
}
