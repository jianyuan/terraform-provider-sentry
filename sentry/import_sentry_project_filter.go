package sentry

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSentryFilterImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	addrID := d.Id()

	tflog.Debug(ctx, "Importing Sentry filters", map[string]interface{}{"filterID": addrID})

	parts := strings.Split(addrID, "/")

	if len(parts) != 3 {
		return nil, errors.New("project import requires an addr_id of the following schema org-slug/project-slug/filter-id")
	}

	d.Set("organization", parts[0])
	d.Set("project", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
