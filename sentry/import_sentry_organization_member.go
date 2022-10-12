package sentry

import (
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func resourceOrganizationMemberImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	addrID := d.Id()

	tflog.Debug(ctx, "Importing Sentry organization member", map[string]interface{}{
		"addrId": addrID,
	})

	parts := strings.Split(addrID, "/")

	if len(parts) != 2 {
		return nil, errors.New("organization member import requires an ADDR ID of the following schema org-slug/member-id")
	}

	d.Set("organization", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
