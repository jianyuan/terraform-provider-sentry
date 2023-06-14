package sentry

import (
	"context"
	"sort"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func resourceSentryReleaseDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSentryReleaseDeploymentCreate,
		ReadContext:   resourceSentryReleaseDeploymentRead,
		UpdateContext: resourceSentryReleaseDeploymentUpdate,
		DeleteContext: resourceSentryReleaseDeploymentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the deploy belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"version": {
				Description: "The version identifier of the release.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The environment this deployment is for.",
			},
			"url": {
				Description: "The optional URL that points to the deploy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Description: "The optional name of the deploy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"projects": {
				Description: "The optional list of projects to deploy.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceSentryReleaseDeploymentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	version := d.Get("version").(string)
	params := releaseDeploymentCreateParams(d)

	deploy, _, err := client.ReleaseDeployments.Create(ctx, org, version, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildThreePartID(org, version, deploy.ID))
	return resourceSentryReleaseDeploymentRead(ctx, d, meta)
}

func resourceSentryReleaseDeploymentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org, release, deployID, err := splitSentryReleaseDeploymentID(d.Id())

	tflog.Debug(ctx, "Reading release deployment", map[string]interface{}{
		"org":      org,
		"release":  release,
		"deployID": deployID,
	})
	deploy, resp, err := client.ReleaseDeployments.Get(ctx, org, release, deployID)
	if found, err := checkClientGet(resp, err, d); !found {
		tflog.Info(ctx, "Removed deployment from state because it no longer exists in Sentry", map[string]interface{}{
			"org":      org,
			"release":  release,
			"deployID": deployID,
		})
		return diag.FromErr(err)
	}

	sort.Strings(deploy.Projects)

	d.SetId(buildThreePartID(org, release, deploy.ID))
	retErr := multierror.Append(
		d.Set("organization", org),
		d.Set("version", release),
		d.Set("environment", deploy.Environment),
		d.Set("url", deploy.URL),
		d.Set("name", deploy.Name),
		d.Set("projects", deploy.Projects),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}

func resourceSentryReleaseDeploymentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	// Since we cannot update or delete deployments we create a new deployment
	// when an update is needed
	org := d.Get("organization").(string)
	version := d.Get("version").(string)
	params := releaseDeploymentCreateParams(d)

	deploy, _, err := client.ReleaseDeployments.Create(ctx, org, version, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildThreePartID(org, version, deploy.ID))
	return resourceSentryReleaseDeploymentRead(ctx, d, meta)
}

func resourceSentryReleaseDeploymentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	org, release, deployID, err := splitSentryReleaseDeploymentID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Deleting deployment", map[string]interface{}{
		"org":      org,
		"release":  release,
		"deployID": deployID,
	})

	// Sentry has no option to delete a deployment. So we skip this and just
	// remove it from the terraform state
	return nil
}

func splitSentryReleaseDeploymentID(id string) (org, release, deployID string, err error) {
	org, release, deployID, err = splitThreePartID(id, "organization-id", "release", "deploy-id")
	return
}

func releaseDeploymentCreateParams(d *schema.ResourceData) *sentry.ReleaseDeployment {
	params := &sentry.ReleaseDeployment{
		Environment: d.Get("environment").(string),
	}
	if v := d.Get("name").(string); v != "" {
		params.Name = sentry.String(v)
	}
	if v := d.Get("url").(string); v != "" {
		params.URL = sentry.String(v)
	}
	if v, ok := d.GetOk("projects"); ok {
		projects := expandStringList(v.([]interface{}))
		if len(projects) > 0 {
			params.Projects = projects
		}
	}
	return params

}
