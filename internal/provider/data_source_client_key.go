package provider

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ datasource.DataSource = &ClientKeyDataSource{}
var _ datasource.DataSourceWithConfigure = &ClientKeyDataSource{}
var _ datasource.DataSourceWithConfigValidators = &ClientKeyDataSource{}

func NewClientKeyDataSource() datasource.DataSource {
	return &ClientKeyDataSource{}
}

type ClientKeyDataSource struct {
	baseDataSource
}

type ClientKeyDataSourceModel struct {
	Organization    types.String `tfsdk:"organization"`
	Project         types.String `tfsdk:"project"`
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	First           types.Bool   `tfsdk:"first"`
	ProjectId       types.String `tfsdk:"project_id"`
	RateLimitWindow types.Int64  `tfsdk:"rate_limit_window"`
	RateLimitCount  types.Int64  `tfsdk:"rate_limit_count"`
	Public          types.String `tfsdk:"public"`
	Secret          types.String `tfsdk:"secret"`
	DsnPublic       types.String `tfsdk:"dsn_public"`
	DsnSecret       types.String `tfsdk:"dsn_secret"`
	DsnCsp          types.String `tfsdk:"dsn_csp"`
}

func (m *ClientKeyDataSourceModel) Fill(organization string, project string, first bool, key sentry.ProjectKey) error {
	m.Organization = types.StringValue(organization)
	m.Project = types.StringValue(project)
	m.Id = types.StringValue(key.ID)
	m.Name = types.StringValue(key.Name)
	m.First = types.BoolValue(first)

	if key.RateLimit == nil {
		m.RateLimitWindow = types.Int64Null()
		m.RateLimitCount = types.Int64Null()
	} else {
		m.RateLimitWindow = types.Int64Value(int64(key.RateLimit.Window))
		m.RateLimitCount = types.Int64Value(int64(key.RateLimit.Count))
	}

	m.ProjectId = types.StringValue(key.ProjectID.String())
	m.Public = types.StringValue(key.Public)
	m.Secret = types.StringValue(key.Secret)
	m.DsnPublic = types.StringValue(key.DSN.Public)
	m.DsnSecret = types.StringValue(key.DSN.Secret)
	m.DsnCsp = types.StringValue(key.DSN.CSP)

	return nil
}

func (d *ClientKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

func (d *ClientKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve a Project's Client Key.",

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization the resource belongs to.",
				Required:            true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The slug of the project the resource belongs to.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the client key.",
				Optional:            true,
			},
			"first": schema.BoolAttribute{
				MarkdownDescription: "Boolean flag indicating that we want the first key of the returned keys.",
				Optional:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project that the key belongs to.",
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
			"rate_limit_window": schema.Int64Attribute{
				MarkdownDescription: "Length of time that will be considered when checking the rate limit.",
				Computed:            true,
			},
			"rate_limit_count": schema.Int64Attribute{
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
	}
}

func (d *ClientKeyDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
			path.MatchRoot("first"),
		),
	}
}

func (d *ClientKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClientKeyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var foundKey *sentry.ProjectKey

	if data.Id.IsNull() {
		var allKeys []*sentry.ProjectKey
		params := &sentry.ListProjectKeysParams{}
		for {
			keys, apiResp, err := d.client.ProjectKeys.List(ctx, data.Organization.ValueString(), data.Project.ValueString(), params)
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", err))
				return
			}

			allKeys = append(allKeys, keys...)

			if apiResp.Cursor == "" {
				break
			}
			params.Cursor = apiResp.Cursor
		}

		if data.Name.IsNull() {
			if len(allKeys) == 1 {
				foundKey = allKeys[0]
			} else if !data.First.IsNull() && data.First.ValueBool() {
				// Find the first key

				// Sort keys by date created
				sort.Slice(allKeys, func(i, j int) bool {
					return allKeys[i].DateCreated.Before(allKeys[j].DateCreated)
				})

				foundKey = allKeys[0]
			} else {

				resp.Diagnostics.AddError("Client Error", "Multiple keys found, please specify the key by `name`, `id`, or set the `first` flag to `true`.")
				return
			}
		} else {
			// Find the key by name
			for _, key := range allKeys {
				if key.Name == data.Name.ValueString() {
					foundKey = key
					break
				}
			}
		}

	} else {
		// Get the key by ID
		key, apiResp, err := d.client.ProjectKeys.Get(
			ctx,
			data.Organization.ValueString(),
			data.Project.ValueString(),
			data.Id.ValueString(),
		)
		if apiResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Not found: %s", err.Error()))
			return
		}
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", err.Error()))
			return
		}

		foundKey = key
	}

	if foundKey == nil {
		resp.Diagnostics.AddError("Client Error", "No key found")
		return
	}

	if err := data.Fill(
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.First.ValueBool(),
		*foundKey,
	); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Fill error: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
