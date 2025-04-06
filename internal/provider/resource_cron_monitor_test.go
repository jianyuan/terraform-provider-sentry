package provider_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCronMonitorResource_basic(t *testing.T) {
	name := randomString()
	resourceName := "sentry_cron_monitor.test"

	r := resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCronMonitorConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "schedule", "* * * * *"),
					resource.TestCheckResourceAttr(resourceName, "timezone", "UTC"),
				),
			},
		},
	})

	if r != nil {
		t.Fatalf("Test failed: %v", r)
	}
}

func testAccCronMonitorConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "sentry_cron_monitor" "test" {
  organization = "%s"
  project      = "%s"
  name         = "%s"
  schedule     = "* * * * *"
  schedule_type = "crontab"
  timezone     = "UTC"
  enabled      = true
}
`, testAccOrganization, testAccProject, name)
}

func randomString() string {
	return strings.ToLower(randomStringFromCharset(10, charset))
}

const charset = "abcdefghijklmnopqrstuvwxyz"

func randomStringFromCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
