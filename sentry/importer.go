package sentry

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func importOrganizationAndID(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected organization-slug/id", id)
	}

	d.Set("organization", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func importOrganizationProjectAndID(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	parts := strings.SplitN(id, "/", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected organization-slug/project-slug/id", id)
	}

	d.Set("organization", parts[0])
	d.Set("project", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
