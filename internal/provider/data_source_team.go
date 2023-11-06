package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ datasource.DataSource = &TeamDataSource{}

func NewTeamDataSource() datasource.DataSource {
	return &TeamDataSource{}
}

type TeamDataSource struct {
	client *sentry.Client
}

type TeamDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	InternalId   types.String `tfsdk:"internal_id"`
	HasAccess    types.Bool   `tfsdk:"has_access"`
	IsPending    types.Bool   `tfsdk:"is_pending"`
	IsMember     types.Bool   `tfsdk:"is_member"`
}

func (d *TeamDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (d *TeamDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Team data source.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
			},
			"organization": schema.StringAttribute{
				Description: "The slug of the organization the team should be created for.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the team.",
				Computed:    true,
			},
			"slug": schema.StringAttribute{
				Description: "The optional slug for this team.",
				Required:    true,
			},
			"internal_id": schema.StringAttribute{
				Description: "The internal ID for this team.",
				Computed:    true,
			},
			"has_access": schema.BoolAttribute{
				Description: "Whether the authenticated user has access to this team.",
				Computed:    true,
			},
			"is_pending": schema.BoolAttribute{
				Description: "Whether the team is pending.",
				Computed:    true,
			},
			"is_member": schema.BoolAttribute{
				Description: "Whether the authenticated user is a member of this team.",
				Computed:    true,
			},
		},
	}
}

func (d *TeamDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sentry.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sentry.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *TeamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TeamDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, _, err := d.client.Teams.Get(ctx, data.Organization.ValueString(), data.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read team, got error: %s", err))
		return
	}

	data.Id = types.StringPointerValue(team.Slug)
	data.Name = types.StringPointerValue(team.Name)
	data.Slug = types.StringPointerValue(team.Slug)
	data.InternalId = types.StringPointerValue(team.ID)
	data.HasAccess = types.BoolPointerValue(team.HasAccess)
	data.IsPending = types.BoolPointerValue(team.IsPending)
	data.IsMember = types.BoolPointerValue(team.IsMember)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
