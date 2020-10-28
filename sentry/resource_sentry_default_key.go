package sentry

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
)

func resourceSentryDefaultKey() *schema.Resource {
	// reuse read and update operations
	dKey := resourceSentryKey()
	dKey.Create = resourceSentryDefaultKeyCreate
	dKey.Delete = resourceAwsDefaultVpcDelete

	// Key name is a computed resource for default key
	dKey.Schema["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The name of the key",
	}

	return dKey
}

func resourceSentryDefaultKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	keys, resp, err := client.ProjectKeys.List(org, project)
	if found, err := checkClientGet(resp, err, d); !found {
		return err
	}

	if len(keys) != 1 {
		return fmt.Errorf("Default key not found on the project")
	}

	id := keys[0].ID
	params := &sentry.UpdateProjectKeyParams{
		Name: d.Get("name").(string),
		RateLimit: &sentry.ProjectKeyRateLimit{
			Window: d.Get("rate_limit_window").(int),
			Count:  d.Get("rate_limit_count").(int),
		},
	}

	if _, _, err = client.ProjectKeys.Update(org, project, id, params); err != nil {
		return err
	}

	d.SetId(id)
	return resourceSentryKeyRead(d, meta)
}

func resourceAwsDefaultVpcDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy Default Key. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
