package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
)

type OrganizationMemberDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Email        types.String `tfsdk:"email"`
	Role         types.String `tfsdk:"role"`
}

func (m *OrganizationMemberDataSourceModel) Fill(organization string, d sentry.OrganizationMember) error {
	m.Id = types.StringValue(d.ID)
	m.Organization = types.StringValue(organization)
	m.Email = types.StringValue(d.Email)
	m.Role = types.StringValue(d.OrgRole)

	return nil
}

var _ datasource.DataSource = &OrganizationMemberDataSource{}
var _ datasource.DataSourceWithConfigure = &OrganizationMemberDataSource{}

func NewOrganizationMemberDataSource() datasource.DataSource {
	return &OrganizationMemberDataSource{}
}

type OrganizationMemberDataSource struct {
	baseDataSource
}

func (d *OrganizationMemberDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_member"
}

func (d *OrganizationMemberDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve an organization member by email.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Computed:            true,
			},
			"organization": DataSourceOrganizationAttribute(),
			"email": schema.StringAttribute{
				MarkdownDescription: "The email of the organization member.",
				Required:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "This is the role of the organization member.",
				Computed:            true,
			},
		},
	}
}

func (d *OrganizationMemberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationMemberDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var foundMember *sentry.OrganizationMember
	params := &sentry.ListCursorParams{}

out:
	for {
		members, apiResp, err := d.client.OrganizationMembers.List(ctx, data.Organization.ValueString(), params)
		if err != nil {
			diagutils.AddClientError(resp.Diagnostics, "read", err)
			return
		}

		for _, member := range members {
			if member.Email == data.Email.ValueString() {
				foundMember = member
				break out
			}
		}

		if apiResp.Cursor == "" {
			break
		}
		params.Cursor = apiResp.Cursor
	}

	if foundMember == nil {
		resp.Diagnostics.AddError("Not found", "No matching organization member found")
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *foundMember); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
