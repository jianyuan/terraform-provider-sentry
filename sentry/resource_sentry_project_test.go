package sentry

import (
	"fmt"
)

func testAccSentryProjectConfig_team(teamName, projectName string) string {
	return testAccSentryTeamConfig(teamName) + fmt.Sprintf(`
resource "sentry_project" "test" {
	organization = sentry_team.test.organization
	teams        = [sentry_team.test.slug]
	name         = "%[1]s"
	platform     = "go"
}
	`, projectName)
}
