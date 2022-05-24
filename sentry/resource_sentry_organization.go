package sentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/sentry"
)

func resourceSentryOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSentryOrganizationCreate,
		ReadContext:   resourceSentryOrganizationRead,
		UpdateContext: resourceSentryOrganizationUpdate,
		DeleteContext: resourceSentryOrganizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The human readable name for the organization",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique URL slug for this organization",
				Computed:    true,
			},
			"agree_terms": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "You agree to the applicable terms of service and privacy policy",
			},
		},
	}
}

func resourceSentryOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	params := &sentry.CreateOrganizationParams{
		Name:       d.Get("name").(string),
		Slug:       d.Get("slug").(string),
		AgreeTerms: sentry.Bool(d.Get("agree_terms").(bool)),
	}

	tflog.Debug(ctx, "Creating Sentry organization", "orgName", params.Name)
	org, _, err := client.Organizations.Create(params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created Sentry organization", "orgName", org.Name, "orgID", org.ID)

	d.SetId(org.Slug)
	return resourceSentryOrganizationRead(ctx, d, meta)
}

func resourceSentryOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Id()

	tflog.Debug(ctx, "Reading Sentry organization", "orgSlug", slug)
	org, resp, err := client.Organizations.Get(slug)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read Sentry organization", "orgSlug", org.Slug, "orgID", org.ID)

	d.SetId(org.Slug)
	d.Set("internal_id", org.ID)
	d.Set("name", org.Name)
	d.Set("slug", org.Slug)
	return nil
}

func resourceSentryOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Id()
	params := &sentry.UpdateOrganizationParams{
		Name: d.Get("name").(string),
		Slug: d.Get("slug").(string),
	}

	tflog.Debug(ctx, "Updating Sentry organization", "orgSlug", slug)
	org, _, err := client.Organizations.Update(slug, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated Sentry organization", "orgSlug", org.Slug, "orgID", org.ID)

	d.SetId(org.Slug)
	return resourceSentryOrganizationRead(ctx, d, meta)
}

func resourceSentryOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Id()

	tflog.Debug(ctx, "Deleting Sentry organization", "orgSlug", slug)
	_, err := client.Organizations.Delete(slug)
	tflog.Debug(ctx, "Deleted Sentry organization", "orgSlug", slug)

	return diag.FromErr(err)
}
