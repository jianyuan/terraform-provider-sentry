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

	check := func(dashboardTitle string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckSentryDashboardExists(rn, &dashboard),
			resource.TestCheckResourceAttr(rn, "organization", testOrganization),
			resource.TestCheckResourceAttr(rn, "title", dashboardTitle),
			resource.TestCheckResourceAttr(rn, "widget.#", "1"),
			resource.TestCheckResourceAttr(rn, "widget.0.title", "Custom Widget"),
			resource.TestCheckResourceAttr(rn, "widget.0.display_type", "world_map"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.#", "1"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.name", "Metric"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.fields.#", "1"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.fields.0", "count()"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.aggregates.#", "1"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.aggregates.0", "count()"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.columns.#", "0"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.field_aliases.#", "0"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.conditions", "!event.type:transaction"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.order_by", ""),
			resource.TestCheckResourceAttrSet(rn, "widget.0.query.0.id"),
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
