package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ datasource.DataSource = &ProjectDataSource{}
var _ datasource.DataSourceWithConfigure = &ProjectDataSource{}

func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

type ProjectDataSource struct {
	baseDataSource
}

type ProjectDataSourceModel struct {
	Organization types.String `tfsdk:"organization"`
	Slug         types.String `tfsdk:"slug"`
	Id           types.String `tfsdk:"id"`
	InternalId   types.String `tfsdk:"internal_id"`
	Name         types.String `tfsdk:"name"`
	Platform     types.String `tfsdk:"platform"`
	DateCreated  types.String `tfsdk:"date_created"`
	Features     types.Set    `tfsdk:"features"`
	Color        types.String `tfsdk:"color"`
	IsPublic     types.Bool   `tfsdk:"is_public"`
}

func (m *ProjectDataSourceModel) Fill(organization string, project sentry.Project) error {
	m.Organization = types.StringValue(organization)
	m.Slug = types.StringValue(project.Slug)
	m.Id = types.StringValue(project.Slug)
	m.InternalId = types.StringValue(project.ID)
	m.Name = types.StringValue(project.Name)
	m.Platform = types.StringValue(project.Platform)
	m.DateCreated = types.StringValue(project.DateCreated.String())

	featureElements := []attr.Value{}
	for _, feature := range project.Features {
		featureElements = append(featureElements, types.StringValue(feature))
	}
	m.Features = types.SetValueMust(types.StringType, featureElements)

	m.Color = types.StringValue(project.Color)
	m.IsPublic = types.BoolValue(project.IsPublic)

	return nil
}

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Return a list of projects available to the authenticated session.",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization the resource belongs to.",
				Required:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "The slug of this project.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The slug of this project.",
				Computed:            true,
			},
			"internal_id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of this project.",
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
			"is_public": schema.BoolAttribute{
				MarkdownDescription: "Whether this project is public.",
				Computed:            true,
			},
		},
	}
}

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, apiResp, err := d.client.Projects.Get(ctx, data.Organization.ValueString(), data.Slug.ValueString())
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Not found: %s", err.Error()))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", err.Error()))
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *project); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Fill error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
