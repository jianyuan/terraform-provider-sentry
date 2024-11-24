package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
)

var _ resource.Resource = &OrganizationRepositoryResource{}
var _ resource.ResourceWithConfigure = &OrganizationRepositoryResource{}
var _ resource.ResourceWithImportState = &OrganizationRepositoryResource{}

func NewOrganizationRepositoryResource() resource.Resource {
	return &OrganizationRepositoryResource{}
}

type OrganizationRepositoryResource struct {
	baseResource
}

func (r *OrganizationRepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_repository"
}

func (r *OrganizationRepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Organization Repository resource. This resource manages Sentry's source code management integrations.",

		Attributes: map[string]schema.Attribute{
			"id":           ResourceIdAttribute(),
			"organization": ResourceOrganizationAttribute(),
			"integration_type": schema.StringAttribute{
				MarkdownDescription: "The type of the organization integration. Supported values are `github`, `github_enterprise`, `gitlab`, `vsts` (Azure DevOps), `bitbucket`, and `bitbucket_server`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("github", "github_enterprise", "gitlab", "vsts", "bitbucket", "bitbucket_server"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"integration_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the organization integration. Source from the URL `https://<organization>.sentry.io/settings/integrations/<integration-type>/<integration-id>/` or use the `sentry_organization_integration` data source.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"identifier": schema.StringAttribute{
				MarkdownDescription: "The identifier of the repository. For GitHub, GitLab and BitBucket, it is `{organization}/{repository}`. For VSTS, it is the [repository ID](https://learn.microsoft.com/en-us/rest/api/azure/devops/git/repositories/get#get-a-repository-by-repositoryid).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *OrganizationRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OrganizationRepositoryModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	repo, _, err := r.client.OrganizationRepositories.Create(
		ctx,
		data.Organization.ValueString(),
		sentry.CreateOrganizationRepositoryParams{
			"provider":     "integrations:" + data.IntegrationType.ValueString(),
			"installation": data.IntegrationId.ValueString(),
			"identifier":   data.Identifier.ValueString(),
		},
	)
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "create", err)
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *repo); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OrganizationRepositoryModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var foundRepo *sentry.OrganizationRepository
	params := &sentry.ListOrganizationRepositoriesParams{
		IntegrationId: data.IntegrationId.ValueString(),
	}

out:
	for {
		repos, apiResp, err := r.client.OrganizationRepositories.List(ctx, data.Organization.ValueString(), params)
		if err != nil {
			diagutils.AddClientError(resp.Diagnostics, "read", err)
			return
		}

		for _, repo := range repos {
			if repo.ID == data.Id.ValueString() {
				foundRepo = repo
				break out
			}
		}

		if apiResp.Cursor == "" {
			break
		}
		params.Cursor = apiResp.Cursor
	}

	if foundRepo == nil {
		diagutils.AddNotFoundError(resp.Diagnostics, "organization repository")
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *foundRepo); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	diagutils.AddNotSupportedError(resp.Diagnostics, "update")
}

func (r *OrganizationRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OrganizationRepositoryModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, apiResp, err := r.client.OrganizationRepositories.Delete(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "delete", err)
		return
	}
}

func (r *OrganizationRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, integrationType, integrationId, id, err := splitFourPartID(req.ID, "organization", "integration_type", "integration_id", "id")
	if err != nil {
		diagutils.AddImportError(resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), organization,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("integration_type"), integrationType,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("integration_id"), integrationId,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), id,
	)...)
}
