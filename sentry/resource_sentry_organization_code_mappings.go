package sentry

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/providerdata"
	"github.com/jianyuan/terraform-provider-sentry/internal/resourceid"
)

const organizationCodeMappingsPerPage = 100

// organizationCodeMappingsListPath builds GET /code-mappings/ with an optional project
// filter. Sentry's unfiltered list uses unordered OFFSET pagination, so orgs with
// more than one page of mappings can skip IDs. Filtering by project keeps each
// page small and stable.
func organizationCodeMappingsListPath(org, integrationId, projectId, cursor string) string {
	q := url.Values{}
	q.Set("per_page", fmt.Sprintf("%d", organizationCodeMappingsPerPage))
	if integrationId != "" {
		q.Set("integrationId", integrationId)
	}
	if projectId != "" {
		q.Set("project", projectId)
	}
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	return fmt.Sprintf("0/organizations/%s/code-mappings/?%s", org, q.Encode())
}

func listOrganizationCodeMappings(ctx context.Context, client *sentry.Client, org, integrationId, projectId string) ([]*sentry.OrganizationCodeMapping, error) {
	var all []*sentry.OrganizationCodeMapping
	seen := make(map[string]struct{})
	cursor := ""
	for page := 0; page < 50; page++ {
		req, err := client.NewRequest("GET", organizationCodeMappingsListPath(org, integrationId, projectId, cursor), nil)
		if err != nil {
			return nil, err
		}
		var pageItems []*sentry.OrganizationCodeMapping
		resp, err := client.Do(ctx, req, &pageItems)
		if err != nil {
			return nil, err
		}
		for _, mapping := range pageItems {
			if mapping == nil || mapping.ID == "" {
				continue
			}
			if _, ok := seen[mapping.ID]; ok {
				continue
			}
			seen[mapping.ID] = struct{}{}
			all = append(all, mapping)
		}
		if resp == nil || resp.Cursor == "" {
			break
		}
		cursor = resp.Cursor
	}
	return all, nil
}

func findOrganizationCodeMapping(mappings []*sentry.OrganizationCodeMapping, id string) *sentry.OrganizationCodeMapping {
	for _, mapping := range mappings {
		if mapping != nil && mapping.ID == id {
			return mapping
		}
	}
	return nil
}

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
	client := meta.(*providerdata.ProviderData).Client

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
	if err := d.Set("internal_id", orgCodeMapping.ID); err != nil {
		return diag.FromErr(err)
	}
	return resourceSentryOrganizationCodeMappingRead(ctx, d, meta)
}

func resourceSentryOrganizationCodeMappingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerdata.ProviderData).Client

	id := d.Id()
	org := d.Get("organization").(string)

	// on import, if the integration_id is an empty string, the API still works
	// it just does not filter by integration_id so more iterations may be required
	// but id *should* be unique across integrations
	integrationId := d.Get("integration_id").(string)
	projectId := d.Get("project_id").(string)

	tflog.Debug(ctx, "Reading Sentry Organization Code Mapping", map[string]interface{}{
		"id":         id,
		"org":        org,
		"project_id": projectId,
	})

	orgCodeMappings, err := listOrganizationCodeMappings(ctx, client, org, integrationId, projectId)
	if err != nil {
		return diag.FromErr(err)
	}

	orgCodeMapping := findOrganizationCodeMapping(orgCodeMappings, id)
	if orgCodeMapping == nil {
		if projectId != "" {
			tflog.Info(ctx, "Removing organization code mapping from state because it no longer exists in Sentry", map[string]interface{}{
				"id":         id,
				"project_id": projectId,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Can't find Sentry Organization Code Mapping: %s", id)
	}

	d.SetId(orgCodeMapping.ID)
	err = errors.Join(
		d.Set("internal_id", orgCodeMapping.ID),
		d.Set("integration_id", orgCodeMapping.IntegrationId),
		d.Set("repository_id", orgCodeMapping.RepoId),
		d.Set("project_id", orgCodeMapping.ProjectId),
		d.Set("default_branch", orgCodeMapping.DefaultBranch),
		d.Set("stack_root", orgCodeMapping.StackRoot),
		d.Set("source_root", orgCodeMapping.SourceRoot),
	)
	return diag.FromErr(err)
}

func resourceSentryOrganizationCodeMappingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerdata.ProviderData).Client

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
	if err := d.Set("internal_id", orgCodeMapping.ID); err != nil {
		return diag.FromErr(err)
	}

	return resourceSentryOrganizationCodeMappingRead(ctx, d, meta)
}

func resourceSentryOrganizationCodeMappingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerdata.ProviderData).Client

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
	org, id, err := resourceid.Split2Path(d.Id(), "organization-slug", "id")
	if err != nil {
		return nil, err
	}

	d.SetId(id)
	if err := d.Set("organization", org); err != nil {
		return nil, err
	}

	resourceSentryOrganizationCodeMappingRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
