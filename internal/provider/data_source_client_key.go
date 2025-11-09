package provider

import (
	"context"
	"net/http"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/maputils"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type ClientKeyDataSourceModel struct {
	Organization           types.String                                                               `tfsdk:"organization"`
	Project                types.String                                                               `tfsdk:"project"`
	Id                     types.String                                                               `tfsdk:"id"`
	Name                   types.String                                                               `tfsdk:"name"`
	First                  types.Bool                                                                 `tfsdk:"first"`
	ProjectId              types.String                                                               `tfsdk:"project_id"`
	RateLimitWindow        types.Int64                                                                `tfsdk:"rate_limit_window"`
	RateLimitCount         types.Int64                                                                `tfsdk:"rate_limit_count"`
	JavascriptLoaderScript supertypes.SingleNestedObjectValueOf[ClientKeyJavascriptLoaderScriptModel] `tfsdk:"javascript_loader_script"`
	Public                 types.String                                                               `tfsdk:"public"`
	Secret                 types.String                                                               `tfsdk:"secret"`
	Dsn                    types.Map                                                                  `tfsdk:"dsn"`
	DsnPublic              types.String                                                               `tfsdk:"dsn_public"`
	DsnSecret              types.String                                                               `tfsdk:"dsn_secret"`
	DsnCsp                 types.String                                                               `tfsdk:"dsn_csp"`
}

func (m *ClientKeyDataSourceModel) Fill(ctx context.Context, key apiclient.ProjectKey) (diags diag.Diagnostics) {
	m.Id = types.StringValue(key.Id)
	m.ProjectId = types.StringValue(key.ProjectId.String())
	m.Name = types.StringValue(key.Name)

	if v, err := key.RateLimit.Get(); err == nil {
		m.RateLimitWindow = types.Int64Value(int64(v.Window))
		m.RateLimitCount = types.Int64Value(int64(v.Count))
	} else {
		m.RateLimitWindow = types.Int64Null()
		m.RateLimitCount = types.Int64Null()
	}

	var javascriptLoaderScript ClientKeyJavascriptLoaderScriptModel
	diags.Append(javascriptLoaderScript.Fill(ctx, key)...)
	if diags.HasError() {
		return
	}

	m.JavascriptLoaderScript = supertypes.NewSingleNestedObjectValueOfNull[ClientKeyJavascriptLoaderScriptModel](ctx)
	diags.Append(m.JavascriptLoaderScript.Set(ctx, &javascriptLoaderScript)...)
	if diags.HasError() {
		return
	}

	m.Public = types.StringValue(key.Public)
	m.Secret = types.StringValue(key.Secret)

	m.Dsn = types.MapValueMust(types.StringType, maputils.MapValues(key.Dsn, func(v string) attr.Value {
		return types.StringValue(v)
	}))

	if v, ok := key.Dsn["public"]; ok {
		m.DsnPublic = types.StringValue(v)
	} else {
		m.DsnPublic = types.StringNull()
	}

	if v, ok := key.Dsn["secret"]; ok {
		m.DsnSecret = types.StringValue(v)
	} else {
		m.DsnSecret = types.StringNull()
	}

	if v, ok := key.Dsn["csp"]; ok {
		m.DsnCsp = types.StringValue(v)
	} else {
		m.DsnCsp = types.StringNull()
	}

	return nil
}

var _ datasource.DataSource = &ClientKeyDataSource{}
var _ datasource.DataSourceWithConfigure = &ClientKeyDataSource{}
var _ datasource.DataSourceWithConfigValidators = &ClientKeyDataSource{}

func NewClientKeyDataSource() datasource.DataSource {
	return &ClientKeyDataSource{}
}

type ClientKeyDataSource struct {
	baseDataSource
}

func (d *ClientKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

func (d *ClientKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = clientKeySchema().GetDataSource(ctx)
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

	var foundKey *apiclient.ProjectKey

	if data.Id.IsNull() {
		var allKeys []apiclient.ProjectKey
		params := &apiclient.ListProjectClientKeysParams{}
		for {
			httpResp, err := d.apiClient.ListProjectClientKeysWithResponse(
				ctx,
				data.Organization.ValueString(),
				data.Project.ValueString(),
				params,
			)
			if err != nil {
				resp.Diagnostics.Append(diagutils.NewClientError("read", err))
				return
			}
			if httpResp.StatusCode() != http.StatusOK {
				resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
				return
			}

			allKeys = append(allKeys, *httpResp.JSON200...)

			params.Cursor = sentryclient.ParseNextPaginationCursor(httpResp.HTTPResponse)
			if params.Cursor == nil {
				break
			}
		}

		if data.Name.IsNull() {
			if len(allKeys) == 1 {
				foundKey = ptr.Ptr(allKeys[0])
			} else if !data.First.IsNull() && data.First.ValueBool() {
				// Find the first key

				// Sort keys by date created
				sort.Slice(allKeys, func(i, j int) bool {
					return allKeys[i].DateCreated.Before(allKeys[j].DateCreated)
				})

				foundKey = ptr.Ptr(allKeys[0])
			} else {
				resp.Diagnostics.AddError("Client error", "Multiple keys found, please specify the key by `name`, `id`, or set the `first` flag to `true`.")
				return
			}
		} else {
			// Find the key by name
			for _, key := range allKeys {
				if key.Name == data.Name.ValueString() {
					foundKey = ptr.Ptr(key)
					break
				}
			}
		}

	} else {
		// Get the key by ID
		httpResp, err := d.apiClient.GetProjectClientKeyWithResponse(
			ctx,
			data.Organization.ValueString(),
			data.Project.ValueString(),
			data.Id.ValueString(),
		)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("read", err))
			return
		}
		if httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.Append(diagutils.NewClientError("read", err))
			return
		}

		foundKey = httpResp.JSON200
	}

	if foundKey == nil {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("client key"))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *foundKey)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
