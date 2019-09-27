package sentry

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
)

func resourceSentryTeam() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryTeamCreate,
		Read:   resourceSentryTeamRead,
		Update: resourceSentryTeamUpdate,
		Delete: resourceSentryTeamDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSentryTeamImport,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the team should be created for",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the team",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The optional slug for this team",
				Computed:    true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"has_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_pending": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_member": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceSentryTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	params := &sentry.CreateTeamParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}
	log.Printf("[DEBUG] Creating Sentry team %s (Organization: %s)", params.Name, org)

	team, _, err := client.Teams.Create(org, params)
	if err != nil {
		return err
	}

	d.SetId(team.Slug)
	return resourceSentryTeamRead(d, meta)
}

func resourceSentryTeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)
	log.Printf("[DEBUG] Reading Sentry team %s (Organization: %s)", slug, org)

	team, _, err := client.Teams.Get(org, slug)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(team.Slug)
	d.Set("team_id", team.ID)
	d.Set("name", team.Name)
	d.Set("slug", team.Slug)
	d.Set("organization", org)
	d.Set("has_access", team.HasAccess)
	d.Set("is_pending", team.IsPending)
	d.Set("is_member", team.IsMember)
	return nil
}

func resourceSentryTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)
	params := &sentry.UpdateTeamParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}
	log.Printf("[DEBUG] Updating Sentry team %s (Organization: %s)", slug, org)

	team, _, err := client.Teams.Update(org, slug, params)
	if err != nil {
		return err
	}

	d.SetId(team.Slug)
	return resourceSentryTeamRead(d, meta)
}

func resourceSentryTeamDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sentry.Client)

	slug := d.Id()
	org := d.Get("organization").(string)
	log.Printf("[DEBUG] Deleting Sentry team %s (Organization: %s)", slug, org)

	_, err := client.Teams.Delete(org, slug)
	return err
}
