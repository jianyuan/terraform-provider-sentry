package sentry

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSentryRule() *schema.Resource {
	resource := resourceSentryIssueAlert()
	resource.DeprecationMessage = "Use the `sentry_issue_alert` resource instead."
	resource.Description = "> **WARNING:** This resource is deprecated and will be removed in the next major version. Use the `sentry_issue_alert` resource instead."
	return resource
}
