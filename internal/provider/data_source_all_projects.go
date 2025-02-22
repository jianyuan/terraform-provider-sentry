package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/sliceutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type AllProjectsDataSourceProjectModel struct {
	InternalId  types.String                  `tfsdk:"internal_id"`
	Slug        types.String                  `tfsdk:"slug"`
	Name        types.String                  `tfsdk:"name"`
	Platform    types.String                  `tfsdk:"platform"`
	DateCreated types.String                  `tfsdk:"date_created"`
	Features    supertypes.SetValueOf[string] `tfsdk:"features"`
	Color       types.String                  `tfsdk:"color"`
}

func (m *AllProjectsDataSourceProjectModel) Fill(ctx context.Context, project apiclient.Project) (diags diag.Diagnostics) {
	m.InternalId = types.StringValue(project.Id)
	m.Slug = types.StringValue(project.Slug)
	m.Name = types.StringValue(project.Name)
	m.Platform = types.StringPointerValue(project.Platform)
	m.DateCreated = types.StringValue(project.DateCreated.String())
	m.Features = supertypes.NewSetValueOfSlice(ctx, project.Features)
	m.Color = types.StringValue(project.Color)
	return
}

type AllProjectsDataSourceModel struct {
	Organization types.String                        `tfsdk:"organization"`
	ProjectSlugs supertypes.SetValueOf[string]       `tfsdk:"project_slugs"`
	Projects     []AllProjectsDataSourceProjectModel `tfsdk:"projects"`
}

func (m *AllProjectsDataSourceModel) Fill(ctx context.Context, projects []apiclient.Project) (diags diag.Diagnostics) {
	projectSlugs := sliceutils.Map(func(project apiclient.Project) string {
		return project.Slug
	}, projects)
	m.ProjectSlugs = supertypes.NewSetValueOfSlice(ctx, projectSlugs)

	m.Projects = make([]AllProjectsDataSourceProjectModel, len(projects))
	for i, project := range projects {
		diags.Append(m.Projects[i].Fill(ctx, project)...)
	}
	return
}

var _ datasource.DataSource = &AllProjectsDataSource{}
var _ datasource.DataSourceWithConfigure = &AllProjectsDataSource{}

func NewAllProjectsDataSource() datasource.DataSource {
	return &AllProjectsDataSource{}
}

type AllProjectsDataSource struct {
	baseDataSource
}

func (d *AllProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_all_projects"
}

func (d *AllProjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Return a list of projects available to the authenticated session.",

		Attributes: map[string]schema.Attribute{
			"organization": DataSourceOrganizationAttribute(),
			"project_slugs": schema.SetAttribute{
				MarkdownDescription: "The slugs of the projects.",
				CustomType:          supertypes.NewSetTypeOf[string](ctx),
				Computed:            true,
			},
			"projects": schema.SetNestedAttribute{
				MarkdownDescription: "The list of projects.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"internal_id": schema.StringAttribute{
							MarkdownDescription: "The internal ID of this project.",
							Computed:            true,
						},
						"slug": schema.StringAttribute{
							MarkdownDescription: "The slug of this project.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of this project.",
							Computed:            true,
						},
						"platform": schema.StringAttribute{
							MarkdownDescription: "The platform of this project.",
							Computed:            true,
						},
						"date_created": schema.StringAttribute{
							MarkdownDescription: "The date this project was created.",
							Computed:            true,
						},
						"features": schema.SetAttribute{
							MarkdownDescription: "The features of this project.",
							CustomType:          supertypes.NewSetTypeOf[string](ctx),
							Computed:            true,
						},
						"color": schema.StringAttribute{
							MarkdownDescription: "The color of this project.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AllProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AllProjectsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var allProjects []apiclient.Project
	params := &apiclient.ListOrganizationProjectsParams{}

	for {
		httpResp, err := d.apiClient.ListOrganizationProjectsWithResponse(ctx, data.Organization.ValueString(), params)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("read", err))
			return
		} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
			resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
			return
		}

		allProjects = append(allProjects, *httpResp.JSON200...)

		params.Cursor = sentryclient.ParseNextPaginationCursor(httpResp.HTTPResponse)
		if params.Cursor == nil {
			break
		}
	}

	resp.Diagnostics.Append(data.Fill(ctx, allProjects)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
