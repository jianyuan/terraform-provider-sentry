package sentry

import (
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"
)

func resourceKeyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	log.Printf("[DEBUG] Importing key using ADDR ID %s", id)

	parts := strings.Split(id, "/")

	if len(parts) != 3 {
		return nil, errors.New("Key import requires an ADDR ID of the following schema org-slug/project-slug/key-id")
	}

	d.Set("organization", parts[0])
	d.Set("project", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
