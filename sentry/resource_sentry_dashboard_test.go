package sentry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSentryDashboard_basic(t *testing.T) {
	dashboardTitle := acctest.RandomWithPrefix("tf-dashboard")
	rn := "sentry_dashboard.test"

	var dashboardID string

	check := func(dashboardTitle string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckSentryDashboardExists(rn, &dashboardID),
			resource.TestCheckResourceAttr(rn, "organization", testOrganization),
			resource.TestCheckResourceAttr(rn, "title", dashboardTitle),
			resource.TestCheckResourceAttr(rn, "widget.#", "1"),
			resource.TestCheckResourceAttr(rn, "widget.0.title", "Custom Widget"),
			resource.TestCheckResourceAttr(rn, "widget.0.display_type", "table"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.#", "1"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.name", "Metric"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.fields.#", "3"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.fields.0", "geo.country_code"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.fields.1", "geo.region"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.fields.2", "count()"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.aggregates.#", "1"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.aggregates.0", "count()"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.columns.#", "0"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.field_aliases.#", "0"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.conditions", "!event.type:transaction has:geo.country_code"),
			resource.TestCheckResourceAttr(rn, "widget.0.query.0.order_by", ""),
			resource.TestCheckResourceAttrSet(rn, "widget.0.query.0.id"),
			resource.TestCheckResourceAttrPtr(rn, "internal_id", &dashboardID),
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

func testAccCheckSentryDashboardExists(n string, dashboardID *string) resource.TestCheckFunc {
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
		*dashboardID = sentry.StringValue(gotDashboard.ID)
		return nil
	}
}

func testAccSentryDashboardConfig(dashboardTitle string) string {
	return testAccSentryOrganizationDataSourceConfig + fmt.Sprintf(`
resource "sentry_dashboard" "test" {
	organization = data.sentry_organization.test.id
	title        = "%[1]s"

	widget {
		title        = "Custom Widget"
		display_type = "table"

		query {
			name       = "Metric"

			fields     = ["geo.country_code", "geo.region", "count()"]
			aggregates = ["count()"]
			conditions = "!event.type:transaction has:geo.country_code"
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
	`, dashboardTitle)
}
