package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
)

type AllOrganizationMembersDataSourceMemberModel struct {
	Id     types.String `tfsdk:"id"`
	UserId types.String `tfsdk:"user_id"`
	Email  types.String `tfsdk:"email"`
	Role   types.String `tfsdk:"role"`
}

func (m *AllOrganizationMembersDataSourceMemberModel) Fill(ctx context.Context, member *sentry.OrganizationMember) (diags diag.Diagnostics) {
	m.Id = types.StringValue(member.ID)
	if member.User.ID != "" {
		m.UserId = types.StringValue(member.User.ID)
	} else {
		m.UserId = types.StringNull()
	}
	m.Email = types.StringValue(member.Email)
	m.Role = types.StringValue(member.OrgRole)
	return nil
}

type AllOrganizationMembersDataSourceModel struct {
	Organization types.String                                  `tfsdk:"organization"`
	Members      []AllOrganizationMembersDataSourceMemberModel `tfsdk:"members"`
}

func (m *AllOrganizationMembersDataSourceModel) Fill(ctx context.Context, members []*sentry.OrganizationMember) (diags diag.Diagnostics) {
	m.Members = make([]AllOrganizationMembersDataSourceMemberModel, len(members))
	for i, member := range members {
		diags.Append(m.Members[i].Fill(ctx, member)...)
	}
	return
}

var _ datasource.DataSource = &AllOrganizationMembersDataSource{}
var _ datasource.DataSourceWithConfigure = &OrganizationMemberDataSource{}

func NewAllOrganizationMembersDataSource() datasource.DataSource {
	return &AllOrganizationMembersDataSource{}
}

type AllOrganizationMembersDataSource struct {
	baseDataSource
}

func (d *AllOrganizationMembersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_all_organization_members"
}

func (d *AllOrganizationMembersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve all organization members.",

		Attributes: map[string]schema.Attribute{
			"organization": DataSourceOrganizationAttribute(),
			"members": schema.SetNestedAttribute{
				MarkdownDescription: "The list of members.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of of the organization member.",
							Computed:            true,
						},
						"user_id": schema.StringAttribute{
							MarkdownDescription: "The user ID of the organization member.",
							Computed:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "The email of the organization member.",
							Computed:            true,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "This is the role of the organization member.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AllOrganizationMembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AllOrganizationMembersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var allMembers []*sentry.OrganizationMember
	params := &sentry.ListCursorParams{}

	for {
		members, sentryResp, err := d.client.OrganizationMembers.List(ctx, data.Organization.ValueString(), params)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("read", err))
			return
		}

		allMembers = append(allMembers, members...)

		if sentryResp.Cursor == "" {
			break
		}
		params.Cursor = sentryResp.Cursor
	}

	resp.Diagnostics.Append(data.Fill(ctx, allMembers)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
