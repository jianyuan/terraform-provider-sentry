package sentry

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func importOrganizationAndID(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	org, id, err := splitTwoPartID(d.Id(), "organization-slug", "id")
	if err != nil {
		return nil, err
	}

	d.Set("organization", org)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func importOrganizationProjectAndID(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	org, project, id, err := splitThreePartID(d.Id(), "organization-slug", "project-slug", "id")
	if err != nil {
		return nil, err
	}

	d.Set("organization", org)
	d.Set("project", project)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
