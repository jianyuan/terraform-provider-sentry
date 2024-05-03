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
	Id           types.String      `tfsdk:"id"`
	Slug         types.String      `tfsdk:"slug"`
	Name         types.String      `tfsdk:"name"`
	Platform     types.String      `tfsdk:"platform"`
	DateCreated  types.String      `tfsdk:"date_created"`
	Features     types.Set         `tfsdk:"features"`
	Color        types.String      `tfsdk:"color"`
	Status       types.String      `tfsdk:"status"`
	Organization OrganizationModel `tfsdk:"organization"`
}

func (m *AllProjectsDataSourceProjectModel) Fill(project sentry.Project) error {
	m.Id = types.StringValue(project.ID)
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
	m.Status = types.StringValue(project.Status)
	m.Organization = OrganizationModel{}
	if err := m.Organization.Fill(project.Organization); err != nil {
		return err
	}

	return nil
}

type AllProjectsDataSourceModel struct {
	Projects []AllProjectsDataSourceProjectModel `tfsdk:"projects"`
}

func (m *AllProjectsDataSourceModel) Fill(projects []sentry.Project) error {
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
			"projects": schema.SetNestedAttribute{
				MarkdownDescription: "The list of projects.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of this project.",
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
						"status": schema.StringAttribute{
							MarkdownDescription: "The status of this project.",
							Computed:            true,
						},
						"organization": schema.SingleNestedAttribute{
							MarkdownDescription: "The organization associated with this project.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "The ID of this organization.",
									Computed:            true,
								},
								"slug": schema.StringAttribute{
									MarkdownDescription: "The slug of this organization.",
									Computed:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "The name of this organization.",
									Computed:            true,
								},
							},
						},
					},
				},
				Computed: true,
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
	params := &sentry.ListCursorParams{}

	for {
		projects, apiResp, err := d.client.Projects.List(ctx, params)
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

	if err := data.Fill(allProjects); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Fill error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
