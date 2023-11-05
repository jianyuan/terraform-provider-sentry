package sentry

import (
	"fmt"

	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

var testAccSentryOrganizationDataSourceConfig = fmt.Sprintf(`
data "sentry_organization" "test" {
	slug = "%s"
}
`, acctest.TestOrganization)
