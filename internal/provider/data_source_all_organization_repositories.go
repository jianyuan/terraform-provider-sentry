package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

type AllOrganizationRepositoriesDataSourceRepositoryModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Url             types.String `tfsdk:"url"`
	Status          types.String `tfsdk:"status"`
	DateCreated     types.String `tfsdk:"date_created"`
	IntegrationId   types.String `tfsdk:"integration_id"`
	IntegrationType types.String `tfsdk:"integration_type"`
	Identifier      types.String `tfsdk:"identifier"`
	ExternalId      types.String `tfsdk:"external_id"`
	ProviderId      types.String `tfsdk:"provider_id"`
	ProviderName    types.String `tfsdk:"provider_name"`
	ImportId        types.String `tfsdk:"import_id"`
}

func (m *AllOrganizationRepositoriesDataSourceRepositoryModel) Fill(organization string, repo *sentry.OrganizationRepository) (diags diag.Diagnostics) {
	m.Id = types.StringValue(repo.ID)
	m.Name = stringValueOrNull(repo.Name)
	m.Url = stringValueOrNull(repo.Url)
	m.Status = stringValueOrNull(repo.Status)
	m.ExternalId = stringValueOrNull(repo.ExternalId)

	if !repo.DateCreated.IsZero() {
		m.DateCreated = types.StringValue(repo.DateCreated.Format(time.RFC3339))
	} else {
		m.DateCreated = types.StringNull()
	}

	if repo.Provider.ID != "" {
		m.ProviderId = types.StringValue(repo.Provider.ID)
	} else {
		m.ProviderId = types.StringNull()
	}

	if repo.Provider.Name != "" {
		m.ProviderName = types.StringValue(repo.Provider.Name)
	} else {
		m.ProviderName = types.StringNull()
	}

	integrationType := strings.TrimPrefix(repo.Provider.ID, "integrations:")
	if integrationType != "" {
		m.IntegrationType = types.StringValue(integrationType)
	} else {
		m.IntegrationType = types.StringNull()
	}

	if repo.IntegrationId != "" {
		m.IntegrationId = types.StringValue(repo.IntegrationId)
	} else {
		m.IntegrationId = types.StringNull()
	}

	identifier, err := repositoryIdentifierFromExternalSlug(repo.ExternalSlug)
	if err != nil {
		diags.Append(diagutils.NewFillError(err))
		m.Identifier = types.StringNull()
	} else {
		m.Identifier = identifier
	}

	if integrationType != "" && repo.IntegrationId != "" && repo.ID != "" {
		m.ImportId = types.StringValue(tfutils.BuildFourPartId(organization, integrationType, repo.IntegrationId, repo.ID))
	} else {
		m.ImportId = types.StringNull()
	}

	return
}

type AllOrganizationRepositoriesDataSourceModel struct {
	Organization types.String                                           `tfsdk:"organization"`
	Repositories []AllOrganizationRepositoriesDataSourceRepositoryModel `tfsdk:"repositories"`
}

func (m *AllOrganizationRepositoriesDataSourceModel) Fill(organization string, repos []*sentry.OrganizationRepository) (diags diag.Diagnostics) {
	m.Repositories = make([]AllOrganizationRepositoriesDataSourceRepositoryModel, len(repos))
	for i, repo := range repos {
		diags.Append(m.Repositories[i].Fill(organization, repo)...)
	}
	return
}

var _ datasource.DataSource = &AllOrganizationRepositoriesDataSource{}
var _ datasource.DataSourceWithConfigure = &AllOrganizationRepositoriesDataSource{}

func NewAllOrganizationRepositoriesDataSource() datasource.DataSource {
	return &AllOrganizationRepositoriesDataSource{}
}

type AllOrganizationRepositoriesDataSource struct {
	baseDataSource
}

func (d *AllOrganizationRepositoriesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_all_organization_repositories"
}

func (d *AllOrganizationRepositoriesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve all organization repositories. See the [Sentry documentation](https://docs.sentry.io/api/organizations/list-an-organizations-repositories/) for more information.",

		Attributes: map[string]schema.Attribute{
			"organization": DataSourceOrganizationAttribute(),
			"repositories": schema.SetNestedAttribute{
				MarkdownDescription: "The list of repositories.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the repository.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The repository name.",
							Computed:            true,
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "The repository URL.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "The repository status.",
							Computed:            true,
						},
						"date_created": schema.StringAttribute{
							MarkdownDescription: "The repository creation timestamp.",
							Computed:            true,
						},
						"integration_id": schema.StringAttribute{
							MarkdownDescription: "The integration ID backing this repository.",
							Computed:            true,
						},
						"integration_type": schema.StringAttribute{
							MarkdownDescription: "The integration type (for example: `github`).",
							Computed:            true,
						},
						"identifier": schema.StringAttribute{
							MarkdownDescription: "The repository identifier used by the integration.",
							Computed:            true,
						},
						"external_id": schema.StringAttribute{
							MarkdownDescription: "The external repository ID.",
							Computed:            true,
						},
						"provider_id": schema.StringAttribute{
							MarkdownDescription: "The provider ID, such as `integrations:github`.",
							Computed:            true,
						},
						"provider_name": schema.StringAttribute{
							MarkdownDescription: "The provider display name.",
							Computed:            true,
						},
						"import_id": schema.StringAttribute{
							MarkdownDescription: "Convenience import ID for `sentry_organization_repository` (`organization/integration_type/integration_id/id`).",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AllOrganizationRepositoriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AllOrganizationRepositoriesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var allRepos []*sentry.OrganizationRepository
	params := &sentry.ListOrganizationRepositoriesParams{
		Status: "",
	}

	for {
		repos, sentryResp, err := d.client.OrganizationRepositories.List(ctx, data.Organization.ValueString(), params)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("read", err))
			return
		}

		allRepos = append(allRepos, repos...)

		if sentryResp.Cursor == "" {
			break
		}
		params.Cursor = sentryResp.Cursor
	}

	resp.Diagnostics.Append(data.Fill(data.Organization.ValueString(), allRepos)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func repositoryIdentifierFromExternalSlug(raw json.RawMessage) (types.String, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return types.StringNull(), nil
	}

	var identifierStr string
	if err := json.Unmarshal(raw, &identifierStr); err == nil {
		return stringValueOrNull(identifierStr), nil
	}

	var identifierNum json.Number
	if err := json.Unmarshal(raw, &identifierNum); err == nil {
		return types.StringValue(identifierNum.String()), nil
	}

	return types.StringNull(), fmt.Errorf("failed to unmarshal repository identifier: %s", raw)
}

func stringValueOrNull(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}
