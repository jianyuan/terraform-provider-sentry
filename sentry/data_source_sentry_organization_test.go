package sentry

import (
	"fmt"
)

var testAccSentryOrganizationDataSourceConfig = fmt.Sprintf(`
data "sentry_organization" "test" {
	slug = "%s"
}
`, testOrganization)
