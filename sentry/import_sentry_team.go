package sentry

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

func resourceSentryTeamImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	addrID := d.Id()

	logging.Debugf("Importing key using ADDR ID %s", addrID)

	parts := strings.Split(addrID, "/")

	if len(parts) != 2 {
		return nil, errors.New("Project import requires an ADDR ID of the following schema org-slug/team-slug")
	}

	d.Set("organization", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
