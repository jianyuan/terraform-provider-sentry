package sentry

import (
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSentryRuleImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	addrID := d.Id()

	log.Printf("[DEBUG] Importing rule using ADDR ID %s", addrID)

	parts := strings.Split(addrID, "/")

	if len(parts) != 3 {
		return nil, errors.New("Project import requires an ADDR ID of the following schema org-slug/project-slug/rule-id")
	}

	d.Set("organization", parts[0])
	d.Set("project", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
