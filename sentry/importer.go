package sentry

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func importOrganizationAndID(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	org, id, err := splitTwoPartID(d.Id(), "organization-slug", "id")
	if err != nil {
		return nil, err
	}

	retErr := multierror.Append(d.Set("organization", org))
	d.SetId(id)

	if err := retErr.ErrorOrNil(); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func importOrganizationProjectAndID(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	org, project, id, err := splitThreePartID(d.Id(), "organization-slug", "project-slug", "id")
	if err != nil {
		return nil, err
	}

	retErr := multierror.Append(
		d.Set("organization", org),
		d.Set("project", project),
	)
	d.SetId(id)

	if err := retErr.ErrorOrNil(); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
