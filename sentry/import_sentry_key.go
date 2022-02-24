package sentry

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKeyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	tflog.Debug(ctx, "Importing Sentry key", "keyID", id)

	parts := strings.Split(id, "/")

	if len(parts) != 3 {
		return nil, errors.New("Key import requires an ADDR ID of the following schema org-slug/project-slug/key-id")
	}

	d.Set("organization", parts[0])
	d.Set("project", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
