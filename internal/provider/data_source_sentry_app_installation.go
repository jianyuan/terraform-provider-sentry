package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
)

type SentryAppInstallationDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Slug         types.String `tfsdk:"slug"`
	Uuid         types.String `tfsdk:"uuid"`
	AppUuid      types.String `tfsdk:"app_uuid"`
	SentryAppId  types.Int64  `tfsdk:"sentry_app_id"`
	Status       types.String `tfsdk:"status"`
}

func (m *SentryAppInstallationDataSourceModel) Fill(d apiclient.SentryAppInstallation) {
	m.Id = types.StringValue(d.Uuid)
	m.Uuid = types.StringValue(d.Uuid)
	m.AppUuid = types.StringValue(d.App.Uuid)
	m.SentryAppId = types.Int64Value(int64(d.App.SentryAppId))
	m.Status = types.StringValue(d.Status)
}

var _ datasource.DataSource = &SentryAppInstallationDataSource{}
var _ datasource.DataSourceWithConfigure = &SentryAppInstallationDataSource{}

func NewSentryAppInstallationDataSource() datasource.DataSource {
	return &SentryAppInstallationDataSource{}
}

type SentryAppInstallationDataSource struct {
	baseDataSource
}

func (d *SentryAppInstallationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_installation"
}

func (d *SentryAppInstallationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Look up a Sentry App Installation by app slug.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
			},
			"organization": DataSourceOrganizationAttribute(),
			"slug": schema.StringAttribute{
				Description: "The slug of the Sentry App to look up.",
				Required:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The installation UUID. Use as `sentry_app_installation_uuid` in `notify_event_sentry_app` issue alert actions.",
				Computed:    true,
			},
			"app_uuid": schema.StringAttribute{
				Description: "The Sentry App UUID.",
				Computed:    true,
			},
			"sentry_app_id": schema.Int64Attribute{
				Description: "The numerical Sentry App ID. Use directly as `sentry_app_id` in metric alert `sentry_app` trigger actions; use `tostring(...)` for `target_identifier`.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The installation status.",
				Computed:    true,
			},
		},
	}
}

func (d *SentryAppInstallationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SentryAppInstallationDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var matched []apiclient.SentryAppInstallation
	params := &apiclient.ListSentryAppInstallationsParams{}

	for {
		httpResp, err := d.apiClient.ListSentryAppInstallationsWithResponse(ctx, data.Organization.ValueString(), params)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("read", err))
			return
		} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
			resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
			return
		}

		for _, installation := range *httpResp.JSON200 {
			if installation.App.Slug == data.Slug.ValueString() {
				matched = append(matched, installation)
			}
		}

		params.Cursor = sentryclient.ParseNextPaginationCursor(httpResp.HTTPResponse)
		if params.Cursor == nil {
			break
		}
	}

	if len(matched) == 0 {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("sentry app installation"))
		return
	} else if len(matched) > 1 {
		resp.Diagnostics.AddError("Not unique", "More than one matching sentry app installation found")
		return
	}

	data.Fill(matched[0])

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
