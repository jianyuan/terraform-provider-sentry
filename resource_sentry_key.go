package main

import (
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSentryKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryKeyCreate,
		Read:   resourceSentryKeyRead,
		Update: resourceSentryKeyUpdate,
		Delete: resourceSentryKeyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeyImporter,
		},

		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the key should be created for",
			},
			"project": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the project the key should be created for",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the key",
			},
			"dsn_secret": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dsn_public": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dsn_csp": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSentryKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	params := &CreateKeyParams{
		Name: d.Get("name").(string),
	}

	key, _, err := client.CreateKey(org, project, params)
	if err != nil {
		return err
	}

	d.SetId(key.ID)
	return resourceSentryKeyRead(d, meta)
}

func resourceSentryKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] SentryKeyRead %s, %s, %s", org, project, id)

	key, _, err := client.GetKey(org, project, id)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(key.ID)
	d.Set("name", key.Label)
	d.Set("secret", key.Secret)
	d.Set("dsn_secret", key.DSN.Secret)
	d.Set("dsn_public", key.DSN.Public)
	d.Set("dsn_csp", key.DSN.CSP)
	return nil
}

func resourceSentryKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)
	params := &UpdateKeyParams{
		Name: d.Get("name").(string),
	}

	key, _, err := client.UpdateKey(org, project, id, params)
	if err != nil {
		return err
	}

	d.SetId(key.ID)
	return resourceSentryKeyRead(d, meta)
}

func resourceSentryKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	_, err := client.DeleteKey(org, project, id)
	return err
}

func resourceKeyImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	addrID := d.Id()

	log.Printf("[DEBUG] Importing key using ADDR ID %s", addrID)

	parts := strings.Split(addrID, "/")

	if len(parts) != 3 {
		return nil, errors.New("Key import requires an ADDR ID of the following schema org-slug/project-slug/key-id")
	}

	d.Set("organization", parts[0])
	d.Set("project", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
