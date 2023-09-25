package sentry

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

// no UpdateContext, unsupported by this integration. will have to ForceNew
func resourceSentryOrganizationRepositoryGitlab() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Gitlab Organization Repository resource.",

		CreateContext: resourceSentryOrganizationRepositoryGitlabCreate,
		ReadContext:   resourceSentryOrganizationRepositoryGitlabRead,
		DeleteContext: resourceSentryOrganizationRepositoryGitlabDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importSentryOrganizationRepositoryGitlab,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the Sentry organization this resource belongs to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"integration_id": {
				Description: "The organization integration ID for Gitlab.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"identifier": {
				Description: "The repo identifier. For Gitlab it is {gitlab_group}/{gitlab_repo}.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"internal_id": {
				Description: "The internal ID for this organization repository.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSentryOrganizationRepositoryGitlabCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	integrationId := d.Get("integration_id").(string)
	identifier := d.Get("identifier").(string)

	tflog.Debug(ctx, "Creating Sentry Gitlab Organization Repository", map[string]interface{}{
		"org":           org,
		"integrationId": integrationId,
		"identifier":    identifier,
	})

	provider := "integrations:gitlab"
	params := sentry.CreateOrganizationRepositoryParams{
		"provider":     provider,
		"installation": integrationId,
		"identifier":   identifier,
	}
	orgRepo, _, err := client.OrganizationRepositories.Create(ctx, org, params)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, "Created Sentry Gitlab Organization Repository", map[string]interface{}{
		"provider":      provider,
		"integrationId": integrationId,
		"identifier":    identifier,
	})

	// identifier contains gitlab group, which is unique globally across sentry
	// You can connect multiple GitLab groups to one Sentry organization, but you cannot connect a single GitLab group to multiple Sentry organizations.
	// https://docs.sentry.io/product/integrations/source-code-mgmt/gitlab/
	d.SetId(identifier)
	d.Set("internal_id", orgRepo.ID)

	return resourceSentryOrganizationRepositoryGitlabRead(ctx, d, meta)
}

func resourceSentryOrganizationRepositoryGitlabRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)

	tflog.Debug(ctx, "Reading Sentry Gitlab Organization Repository", map[string]interface{}{
		"org": org,
		"id":  id,
	})

	// get all paginated organization repositories with the query
	// query does a fuzzy match on name
	var orgRepos []*sentry.OrganizationRepository
	params := &sentry.ListOrganizationRepositoriesParams{
		ListCursorParams: sentry.ListCursorParams{},
		Query:            id,
	}
	for {
		keys, resp, err := client.OrganizationRepositories.List(ctx, org, params)
		if err != nil {
			return diag.FromErr(err)
		}
		orgRepos = append(orgRepos, keys...)

		tflog.Debug(ctx, "Requested organization repositories list cursor", map[string]interface{}{"cursor": resp.Cursor})
		if resp.Cursor == "" {
			break
		}
		params.ListCursorParams.Cursor = resp.Cursor
	}

	tflog.Debug(ctx, "Reading Sentry Gitlab Organization Repository", map[string]interface{}{
		"org": org,
		"id":  id,
	})

	// filter for first exactly matching name
	for _, orgRepo := range orgRepos {
		if orgRepo.Name == id {
			d.SetId(orgRepo.Name)
			retErr := multierror.Append(
				d.Set("internal_id", orgRepo.ID),
				d.Set("integration_id", orgRepo.IntegrationId),
			)
			return diag.FromErr(retErr.ErrorOrNil())
		}
	}

	return diag.Errorf("Can't find Sentry Organization Repository: %s", id)
}

func resourceSentryOrganizationRepositoryGitlabDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	id := d.Id()
	org := d.Get("organization").(string)
	internalId := d.Get("internal_id").(string)

	tflog.Debug(ctx, "Deleting Sentry Gitlab Organization Repository", map[string]interface{}{
		"org":        org,
		"id":         id,
		"internalId": internalId,
	})
	_, _, err := client.OrganizationRepositories.Delete(ctx, org, internalId)
	tflog.Debug(ctx, "Deleted Sentry Gitlab Organization Repository", map[string]interface{}{
		"org":        org,
		"id":         id,
		"internalId": internalId,
	})

	return diag.FromErr(err)
}

func importSentryOrganizationRepositoryGitlab(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	org, id, err := splitTwoPartID(d.Id(), "organization-slug", "id")
	if err != nil {
		return nil, err
	}

	d.SetId(id)
	d.Set("identifier", id)
	d.Set("organization", org)

	resourceSentryOrganizationRepositoryGitlabRead(ctx, d, meta)

	return []*schema.ResourceData{d}, nil
}
