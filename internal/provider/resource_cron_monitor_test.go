package provider

import (
    "fmt"
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "github.com/hashicorp/terraform-plugin-testing/terraform"
    "github.com/jianyuan/terraform-provider-sentry/internal/acctest"
    "github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

func TestAccCronMonitorResource(t *testing.T) {
    rn := "sentry_cron_monitor.test"
    team := acctest.RandomWithPrefix("tf-team")
    project := acctest.RandomWithPrefix("tf-project")
    name := acctest.RandomWithPrefix("tf-cron-monitor")

    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { acctest.PreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccCronMonitorConfig(team, project, name, "* * * * *", "UTC", true),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
                    resource.TestCheckResourceAttr(rn, "project", project),
                    resource.TestCheckResourceAttr(rn, "name", name),
                    resource.TestCheckResourceAttr(rn, "schedule", "* * * * *"),
                    resource.TestCheckResourceAttr(rn, "schedule_type", "crontab"),
                    resource.TestCheckResourceAttr(rn, "timezone", "UTC"),
                    resource.TestCheckResourceAttr(rn, "enabled", "true"),
                ),
            },
            {
                Config: testAccCronMonitorConfig(team, project, name, "*/5 * * * *", "America/New_York", false),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr(rn, "schedule", "*/5 * * * *"),
                    resource.TestCheckResourceAttr(rn, "timezone", "America/New_York"),
                    resource.TestCheckResourceAttr(rn, "enabled", "false"),
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
                    org := rs.Primary.Attributes["organization"]
                    project := rs.Primary.Attributes["project"]
                    id := rs.Primary.ID
                    return tfutils.BuildThreePartId(org, project, id), nil
                },
                ImportStateVerify: true,
            },
        },
    })
}

func testAccCronMonitorConfig(team, project, name, schedule, timezone string, enabled bool) string {
    return fmt.Sprintf(`
resource "sentry_team" "test" {
    organization = "%s"
    name         = "%s"
    slug         = "%s"
}

resource "sentry_project" "test" {
    organization = sentry_team.test.organization
    teams        = [sentry_team.test.id]
    name         = "%s"
    platform     = "go"
}

resource "sentry_cron_monitor" "test" {
    organization = sentry_team.test.organization
    project      = sentry_project.test.slug
    name         = "%s"
    schedule     = "%s"
    schedule_type = "crontab"
    timezone     = "%s"
    enabled      = %t
}
`, acctest.TestOrganization, team, team, project, name, schedule, timezone, enabled)
}