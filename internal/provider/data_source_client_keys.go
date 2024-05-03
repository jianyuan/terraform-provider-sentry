package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ datasource.DataSource = &ClientKeysDataSource{}

func NewClientKeysDataSource() datasource.DataSource {
	return &ClientKeysDataSource{}
}

type ClientKeysDataSource struct {
	client *sentry.Client
}

type ClientKeysDataSourceKeyModel struct {
	Id              types.String `tfsdk:"id"`
	Organization    types.String `tfsdk:"organization"`
	Project         types.String `tfsdk:"project"`
	ProjectId       types.String `tfsdk:"project_id"`
	Name            types.String `tfsdk:"name"`
	Public          types.String `tfsdk:"public"`
	Secret          types.String `tfsdk:"secret"`
	RateLimitWindow types.Int64  `tfsdk:"rate_limit_window"`
	RateLimitCount  types.Int64  `tfsdk:"rate_limit_count"`
	DsnPublic       types.String `tfsdk:"dsn_public"`
	DsnSecret       types.String `tfsdk:"dsn_secret"`
	DsnCsp          types.String `tfsdk:"dsn_csp"`
}

func (m *ClientKeysDataSourceKeyModel) Fill(organization string, project string, d sentry.ProjectKey) error {
	m.Id = types.StringValue(d.ID)
	m.Organization = types.StringValue(organization)
	m.Project = types.StringValue(project)
	m.ProjectId = types.StringValue(d.ProjectID.String())
	m.Name = types.StringValue(d.Name)
	m.Public = types.StringValue(d.Public)
	m.Secret = types.StringValue(d.Secret)

	if d.RateLimit != nil {
		m.RateLimitWindow = types.Int64Value(int64(d.RateLimit.Window))
		m.RateLimitCount = types.Int64Value(int64(d.RateLimit.Count))
	}

	m.DsnPublic = types.StringValue(d.DSN.Public)
	m.DsnSecret = types.StringValue(d.DSN.Secret)
	m.DsnCsp = types.StringValue(d.DSN.CSP)

	return nil
}

type ClientKeysDataSourceModel struct {
	Organization types.String                   `tfsdk:"organization"`
	Project      types.String                   `tfsdk:"project"`
	FilterStatus types.String                   `tfsdk:"filter_status"`
	Keys         []ClientKeysDataSourceKeyModel `tfsdk:"keys"`
}

func (m *ClientKeysDataSourceModel) Fill(organization string, project string, keys []*sentry.ProjectKey) error {
	m.Organization = types.StringValue(organization)
	m.Project = types.StringValue(project)

	for _, key := range keys {
		var model ClientKeysDataSourceKeyModel
		if err := model.Fill(organization, project, *key); err != nil {
			return err
		}

		m.Keys = append(m.Keys, model)
	}

	return nil
}

func (d *ClientKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keys"
}

func (d *ClientKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List a Project's Client Keys.",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization the resource belongs to.",
				Required:            true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The slug of the project the resource belongs to.",
				Required:            true,
			},
			"filter_status": schema.StringAttribute{
				MarkdownDescription: "Filter client keys by `active` or `inactive`. Defaults to returning all keys if not specified.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"active",
						"inactive",
					),
				},
			},
			"keys": schema.ListNestedAttribute{
				MarkdownDescription: "The list of client keys.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of this resource.",
							Computed:            true,
						},
						"organization": schema.StringAttribute{
							MarkdownDescription: "The slug of the organization the resource belongs to.",
							Computed:            true,
						},
						"project": schema.StringAttribute{
							MarkdownDescription: "The slug of the project the resource belongs to.",
							Computed:            true,
						},
						"project_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the project that the key belongs to.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the client key.",
							Computed:            true,
						},
						"public": schema.StringAttribute{
							MarkdownDescription: "The public key.",
							Computed:            true,
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: "The secret key.",
							Computed:            true,
						},
						"rate_limit_window": schema.NumberAttribute{
							MarkdownDescription: "Length of time that will be considered when checking the rate limit.",
							Computed:            true,
						},
						"rate_limit_count": schema.NumberAttribute{
							MarkdownDescription: "Number of events that can be reported within the rate limit window.",
							Computed:            true,
						},
						"dsn_public": schema.StringAttribute{
							MarkdownDescription: "The DSN tells the SDK where to send the events to.",
							Computed:            true,
						},
						"dsn_secret": schema.StringAttribute{
							MarkdownDescription: "Deprecated DSN includes a secret which is no longer required by newer SDK versions. If you are unsure which to use, follow installation instructions for your language.",
							Computed:            true,
						},
						"dsn_csp": schema.StringAttribute{
							MarkdownDescription: "Security header endpoint for features like CSP and Expect-CT reports.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *ClientKeysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ClientKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClientKeysDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var allKeys []*sentry.ProjectKey
	params := &sentry.ListProjectKeysParams{
		Status: data.FilterStatus.ValueStringPointer(),
	}
	for {
		keys, apiResp, err := d.client.ProjectKeys.List(ctx, data.Organization.ValueString(), data.Project.ValueString(), params)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read client keys, got error: %s", err))
			return
		}

		allKeys = append(allKeys, keys...)

		if apiResp.Cursor == "" {
			break
		}
		params.Cursor = apiResp.Cursor
	}

	if err := data.Fill(data.Organization.ValueString(), data.Project.ValueString(), allKeys); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error filling client keys: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
