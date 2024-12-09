package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/sliceutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
)

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

func (m *ProjectDataSourceModel) Fill(project apiclient.Project) error {
	m.Slug = types.StringValue(project.Slug)
	m.Id = types.StringValue(project.Slug)
	m.InternalId = types.StringValue(project.Id)
	m.Name = types.StringValue(project.Name)
	m.Platform = types.StringPointerValue(project.Platform)
	m.DateCreated = types.StringValue(project.DateCreated.String())
	m.Features = types.SetValueMust(types.StringType, sliceutils.Map(func(v string) attr.Value {
		return types.StringValue(v)
	}, project.Features))
	m.Color = types.StringValue(project.Color)
	m.IsPublic = types.BoolValue(project.IsPublic)

	return nil
}

var _ datasource.DataSource = &ProjectDataSource{}
var _ datasource.DataSourceWithConfigure = &ProjectDataSource{}

func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

type ProjectDataSource struct {
	baseDataSource
}

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Return a list of projects available to the authenticated session.",

		Attributes: map[string]schema.Attribute{
			"organization": DataSourceOrganizationAttribute(),
			"slug":         DataSourceProjectAttribute(),
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

	httpResp, err := d.apiClient.GetOrganizationProjectWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Slug.ValueString(),
	)

	if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("project"))
		return
	}
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	}

	if err := data.Fill(*httpResp.JSON200); err != nil {
		resp.Diagnostics.Append(diagutils.NewFillError(err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
