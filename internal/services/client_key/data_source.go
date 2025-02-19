package client_key

import (
	"context"
	"net/http"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/jianyuan/go-utils/ptr"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/services"
)

var _ datasource.DataSource = &DataSource{}
var _ datasource.DataSourceWithConfigure = &DataSource{}
var _ datasource.DataSourceWithConfigValidators = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	services.BaseDataSource
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = Schema().GetDataSource(ctx)
}

func (d *DataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
			path.MatchRoot("first"),
		),
	}
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var foundKey *apiclient.ProjectKey

	if data.Id.IsNull() {
		var allKeys []apiclient.ProjectKey
		params := &apiclient.ListProjectClientKeysParams{}
		for {
			httpResp, err := d.ApiClient.ListProjectClientKeysWithResponse(
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
		httpResp, err := d.ApiClient.GetProjectClientKeyWithResponse(
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
