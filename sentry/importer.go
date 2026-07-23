package sentry

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/terraform-provider-sentry/internal/resourceid"
)

func importOrganizationAndID(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	org, id, err := resourceid.Split2Path(d.Id(), "organization-slug", "id")
	if err != nil {
		return nil, err
	}

	d.SetId(id)

	err = d.Set("organization", org)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func importOrganizationProjectAndID(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	org, project, id, err := resourceid.Split3Path(d.Id(), "organization-slug", "project-slug", "id")
	if err != nil {
		return nil, err
	}

	err = errors.Join(
		d.Set("organization", org),
		d.Set("project", project),
	)
	if err != nil {
		return nil, err
	}

	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
