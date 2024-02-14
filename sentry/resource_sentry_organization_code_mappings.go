package sentry

import (
	"context"

	"github.com/getkevin/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSentryOrganizationCodeMapping() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Organization Code Mapping resource.",

		CreateContext: resourceSentryOrganizationCodeMappingCreate,
		ReadContext:   resourceSentryOrganizationCodeMappingRead,
		UpdateContext: resourceSentryOrganizationCodeMappingUpdate,
		DeleteContext: resourceSentryOrganizationCodeMappingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importSentryOrganizationCodeMapping,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization the code mapping is under.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"integration_id": {
				Description: "Sentry Organization Integration ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"repository_id": {
				Description: "Sentry Organization Repository ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project_id": {
				Description: "Sentry Project ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"default_branch": {
				Description: "Default branch of your code we fall back to if you do not have commit tracking set up.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"stack_root": {
				Description: "https://docs.sentry.io/product/integrations/source-code-mgmt/github/#stack-trace-linking",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"source_root": {
				Description: "https://docs.sentry.io/product/integrations/source-code-mgmt/github/#stack-trace-linking",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"internal_id": {
				Description: "The internal ID for this resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSentryOrganizationCodeMappingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)

	tflog.Debug(ctx, "Creating Sentry Organization Code Mapping", map[string]interface{}{
		"org": org,
	})

	params := sentry.CreateOrganizationCodeMappingParams{
		IntegrationId: d.Get("integration_id").(string),
		RepositoryId:  d.Get("repository_id").(string),
		ProjectId:     d.Get("project_id").(string),
		DefaultBranch: d.Get("default_branch").(string),
		StackRoot:     d.Get("stack_root").(string),
		SourceRoot:    d.Get("source_root").(string),
	}
	orgCodeMapping, _, err := client.OrganizationCodeMappings.Create(ctx, org, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(orgCodeMapping.ID)
	d.Set("internal_id", orgCodeMapping.ID)
	return resourceSentryOrganizationCodeMappingRead(ctx, d, meta)
}

func resourceSentryOrganizationCodeMappingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)

	// on import, if the integration_id is an empty string, the API still works
	// it just does not filter by integration_id so more iterations may be required
	// but id *should* be unique across integrations
	integrationId := d.Get("integration_id").(string)

	tflog.Debug(ctx, "Reading Sentry Organization Code Mapping", map[string]interface{}{
		"id":  id,
		"org": org,
	})

	// get all paginated organization repositories with the query
	// query does a fuzzy match on name
	var orgCodeMappings []*sentry.OrganizationCodeMapping
	params := &sentry.ListOrganizationCodeMappingsParams{
		ListCursorParams: sentry.ListCursorParams{},
		IntegrationId:    integrationId,
	}
	for {
		keys, resp, err := client.OrganizationCodeMappings.List(ctx, org, params)
		if err != nil {
			return diag.FromErr(err)
		}
		orgCodeMappings = append(orgCodeMappings, keys...)

		tflog.Debug(ctx, "Requested organization code mappings list cursor", map[string]interface{}{"cursor": resp.Cursor})
		if resp.Cursor == "" {
			break
		}
		params.ListCursorParams.Cursor = resp.Cursor
	}

	// filter for first exactly matching name
	for _, orgCodeMapping := range orgCodeMappings {
		if orgCodeMapping.ID == id {
			d.SetId(orgCodeMapping.ID)
			retErr := multierror.Append(
				d.Set("internal_id", orgCodeMapping.ID),
				d.Set("integration_id", orgCodeMapping.IntegrationId),
				d.Set("repository_id", orgCodeMapping.RepoId),
				d.Set("project_id", orgCodeMapping.ProjectId),
				d.Set("default_branch", orgCodeMapping.DefaultBranch),
				d.Set("stack_root", orgCodeMapping.StackRoot),
				d.Set("source_root", orgCodeMapping.SourceRoot),
			)
			return diag.FromErr(retErr.ErrorOrNil())
		}
	}

	return diag.Errorf("Can't find Sentry Organization Code Mapping: %s", id)
}

func resourceSentryOrganizationCodeMappingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	params := sentry.UpdateOrganizationCodeMappingParams{
		IntegrationId: d.Get("integration_id").(string),
		RepositoryId:  d.Get("repository_id").(string),
		ProjectId:     d.Get("project_id").(string),
		DefaultBranch: d.Get("default_branch").(string),
		StackRoot:     d.Get("stack_root").(string),
		SourceRoot:    d.Get("source_root").(string),
	}

	tflog.Debug(ctx, "Updating Sentry Organization Code Mapping", map[string]interface{}{
		"id":  id,
		"org": org,
	})
	orgCodeMapping, _, err := client.OrganizationCodeMappings.Update(ctx, org, id, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(orgCodeMapping.ID)
	d.Set("internal_id", orgCodeMapping.ID)

	return resourceSentryOrganizationCodeMappingRead(ctx, d, meta)
}

func resourceSentryOrganizationCodeMappingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)

	tflog.Debug(ctx, "Deleting Sentry Organization Code Mapping", map[string]interface{}{
		"id":  id,
		"org": org,
	})
	_, err := client.OrganizationCodeMappings.Delete(ctx, org, id)
	return diag.FromErr(err)
}

func importSentryOrganizationCodeMapping(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	org, id, err := splitTwoPartID(d.Id(), "organization-slug", "id")
	if err != nil {
		return nil, err
	}

	d.SetId(id)
	d.Set("organization", org)

	resourceSentryOrganizationCodeMappingRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
