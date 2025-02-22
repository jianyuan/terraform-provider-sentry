package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
)

type OrganizationMemberDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Email        types.String `tfsdk:"email"`
	Role         types.String `tfsdk:"role"`
}

func (m *OrganizationMemberDataSourceModel) Fill(ctx context.Context, member apiclient.OrganizationMember) (diags diag.Diagnostics) {
	m.Id = types.StringValue(member.Id)
	m.Email = types.StringValue(member.Email)
	m.Role = types.StringValue(member.OrgRole)
	return
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

	var foundMember *apiclient.OrganizationMember
	params := &apiclient.ListOrganizationMembersParams{}

out:
	for {
		httpResp, err := d.apiClient.ListOrganizationMembersWithResponse(ctx, data.Organization.ValueString(), params)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("read", err))
			return
		} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
			resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
			return
		}

		for _, member := range *httpResp.JSON200 {
			if member.Email == data.Email.ValueString() {
				foundMember = &member
				break out
			}
		}

		params.Cursor = sentryclient.ParseNextPaginationCursor(httpResp.HTTPResponse)
		if params.Cursor == nil {
			break
		}
	}

	if foundMember == nil {
		resp.Diagnostics.AddError("Not found", "No matching organization member found")
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *foundMember)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
