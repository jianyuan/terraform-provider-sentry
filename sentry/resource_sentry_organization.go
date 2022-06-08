package sentry

import (
	"context"
	"net/http"

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
			"internal_id": {
				Description: "The internal ID for this organization.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSentryOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	params := &sentry.CreateOrganizationParams{
		Name:       sentry.String(d.Get("name").(string)),
		Slug:       sentry.String(d.Get("slug").(string)),
		AgreeTerms: sentry.Bool(d.Get("agree_terms").(bool)),
	}

	tflog.Info(ctx, "Creating organization", map[string]interface{}{"organization": params.Name})
	organization, _, err := client.Organizations.Create(ctx, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*organization.Slug)
	return resourceSentryOrganizationRead(ctx, d, meta)
}

func resourceSentryOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)
	org := d.Id()

	organization, _, err := client.Organizations.Get(ctx, org)
	if err != nil {
		if sErr, ok := err.(*sentry.ErrorResponse); ok {
			if sErr.Response.StatusCode == http.StatusNotFound {
				tflog.Info(ctx, "Removing organization from state because it no longer exists in Sentry", map[string]interface{}{"org": org})
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	d.Set("name", organization.Name)
	d.Set("slug", organization.Slug)
	d.Set("internal_id", organization.ID)
	return nil
}

func resourceSentryOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)
	org := d.Id()
	params := &sentry.UpdateOrganizationParams{
		Name: sentry.String(d.Get("name").(string)),
		Slug: sentry.String(d.Get("slug").(string)),
	}

	tflog.Debug(ctx, "Updating organization", map[string]interface{}{"org": org})
	organization, _, err := client.Organizations.Update(ctx, org, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*organization.Slug)
	return resourceSentryOrganizationRead(ctx, d, meta)
}

func resourceSentryOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)
	org := d.Id()

	tflog.Debug(ctx, "Deleting organization", map[string]interface{}{"org": org})
	_, err := client.Organizations.Delete(ctx, org)
	return diag.FromErr(err)
}
