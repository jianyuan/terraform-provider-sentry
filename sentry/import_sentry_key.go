package sentry

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

func resourceKeyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	logging.Debugf("Importing key using ADDR ID %s", id)

	parts := strings.Split(id, "/")

	if len(parts) != 3 {
		return nil, errors.New("Key import requires an ADDR ID of the following schema org-slug/project-slug/key-id")
	}

	d.Set("organization", parts[0])
	d.Set("project", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
