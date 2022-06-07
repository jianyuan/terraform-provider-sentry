package sentry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func resourceSentryOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Organization resource.",

		CreateContext: resourceSentryOrganizationCreate,
		ReadContext:   resourceSentryOrganizationRead,
		UpdateContext: resourceSentryOrganizationUpdate,
		DeleteContext: resourceSentryOrganizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The human readable name for the organization.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"slug": {
				Description: "The unique URL slug for this organization.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"agree_terms": {
				Description: "You agree to the applicable terms of service and privacy policy.",
				Type:        schema.TypeBool,
				Required:    true,
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

	tflog.Debug(ctx, "Creating Sentry organization", map[string]interface{}{
		"orgName": params.Name,
	})
	org, _, err := client.Organizations.Create(ctx, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Created Sentry organization", map[string]interface{}{
		"orgName": org.Name,
		"orgID":   org.ID,
	})

	d.SetId(org.Slug)
	return resourceSentryOrganizationRead(ctx, d, meta)
}

func resourceSentryOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Id()

	tflog.Debug(ctx, "Reading Sentry organization", map[string]interface{}{
		"orgSlug": slug,
	})
	org, resp, err := client.Organizations.Get(ctx, slug)
	if found, err := checkClientGet(resp, err, d); !found {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Read Sentry organization", map[string]interface{}{
		"orgSlug": org.Slug,
		"orgID":   org.ID,
	})

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

	tflog.Debug(ctx, "Updating Sentry organization", map[string]interface{}{
		"orgSlug": slug,
	})
	org, _, err := client.Organizations.Update(ctx, slug, params)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, "Updated Sentry organization", map[string]interface{}{
		"orgSlug": org.Slug,
		"orgID":   org.ID,
	})

	d.SetId(org.Slug)
	return resourceSentryOrganizationRead(ctx, d, meta)
}

func resourceSentryOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	slug := d.Id()

	tflog.Debug(ctx, "Deleting Sentry organization", map[string]interface{}{
		"orgSlug": slug,
	})
	_, err := client.Organizations.Delete(ctx, slug)
	tflog.Debug(ctx, "Deleted Sentry organization", map[string]interface{}{
		"orgSlug": slug,
	})

	return diag.FromErr(err)
}
