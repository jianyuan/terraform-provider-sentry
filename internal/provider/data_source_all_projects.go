package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ datasource.DataSource = &AllProjectsDataSource{}
var _ datasource.DataSourceWithConfigure = &AllProjectsDataSource{}

func NewAllProjectsDataSource() datasource.DataSource {
	return &AllProjectsDataSource{}
}

type AllProjectsDataSource struct {
	baseDataSource
}

type AllProjectsDataSourceProjectModel struct {
	InternalId  types.String `tfsdk:"internal_id"`
	Slug        types.String `tfsdk:"slug"`
	Name        types.String `tfsdk:"name"`
	Platform    types.String `tfsdk:"platform"`
	DateCreated types.String `tfsdk:"date_created"`
	Features    types.Set    `tfsdk:"features"`
	Color       types.String `tfsdk:"color"`
}

func (m *AllProjectsDataSourceProjectModel) Fill(project sentry.Project) error {
	m.InternalId = types.StringValue(project.ID)
	m.Slug = types.StringValue(project.Slug)
	m.Name = types.StringValue(project.Name)
	m.Platform = types.StringValue(project.Platform)
	m.DateCreated = types.StringValue(project.DateCreated.String())

	featureElements := []attr.Value{}
	for _, feature := range project.Features {
		featureElements = append(featureElements, types.StringValue(feature))
	}
	m.Features = types.SetValueMust(types.StringType, featureElements)

	m.Color = types.StringValue(project.Color)

	return nil
}

type AllProjectsDataSourceModel struct {
	Organization types.String                        `tfsdk:"organization"`
	ProjectSlugs types.Set                           `tfsdk:"project_slugs"`
	Projects     []AllProjectsDataSourceProjectModel `tfsdk:"projects"`
}

func (m *AllProjectsDataSourceModel) Fill(organization string, projects []sentry.Project) error {
	m.Organization = types.StringValue(organization)

	projectSlugElements := []attr.Value{}
	for _, project := range projects {
		projectSlugElements = append(projectSlugElements, types.StringValue(project.Slug))
	}
	m.ProjectSlugs = types.SetValueMust(types.StringType, projectSlugElements)

	for _, project := range projects {
		p := AllProjectsDataSourceProjectModel{}
		if err := p.Fill(project); err != nil {
			return err
		}
		m.Projects = append(m.Projects, p)
	}

	return nil
}

func (d *AllProjectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_all_projects"
}

func (d *AllProjectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Return a list of projects available to the authenticated session.",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization the resource belongs to.",
				Required:            true,
			},
			"project_slugs": schema.SetAttribute{
				MarkdownDescription: "The slugs of the projects.",
				Computed:            true,
				ElementType:         types.StringType,
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
							Computed:            true,
							ElementType:         types.StringType,
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

	var allProjects []sentry.Project
	params := &sentry.ListOrganizationProjectsParams{}

	for {
		projects, apiResp, err := d.client.OrganizationProjects.List(ctx, data.Organization.ValueString(), params)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", err))
			return
		}

		for _, project := range projects {
			allProjects = append(allProjects, *project)
		}

		if apiResp.Cursor == "" {
			break
		}
		params.Cursor = apiResp.Cursor
	}

	if err := data.Fill(data.Organization.ValueString(), allProjects); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Fill error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
