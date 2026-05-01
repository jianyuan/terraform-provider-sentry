package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

var _ datasource.DataSource = &ErrorMonitorDataSource{}

func NewErrorMonitorDataSource() datasource.DataSource {
	return &ErrorMonitorDataSource{}
}

type ErrorMonitorDataSource struct {
	baseDataSource
}

func (d *ErrorMonitorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_error_monitor"
}

func (d *ErrorMonitorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "⚠️ This resource is currently in beta and may be subject to change. It is supported by [New Monitors and Alerts](https://docs.sentry.io/product/new-monitors-and-alerts/) and may not be viewable in the UI today.\n\nRetrieve an Error Monitor by ID. Useful for referencing monitors created outside of Terraform in `sentry_alert.monitor_ids`.",
		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "The organization slug.",
				Required:            true,
				CustomType:          supertypes.StringType{},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of the monitor.",
				Required:            true,
				CustomType:          supertypes.StringType{},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the monitor.",
				Computed:            true,
				CustomType:          supertypes.StringType{},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the monitor is enabled.",
				Computed:            true,
				CustomType:          supertypes.BoolType{},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of the project this monitor belongs to.",
				Computed:            true,
				CustomType:          supertypes.StringType{},
			},
		},
	}
}

func (d *ErrorMonitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ErrorMonitorDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := d.apiClient.GetProjectMonitorWithResponse(ctx, data.Organization.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got error: %s", err))
		return
	} else if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read, got status code %d: %s", httpResp.StatusCode(), string(httpResp.Body)))
		return
	} else if httpResp.JSON200 == nil {
		resp.Diagnostics.AddError("Client Error", "Unable to read, got empty response body")
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type ErrorMonitorDataSourceModel struct {
	Organization supertypes.StringValue `tfsdk:"organization"`
	Id           supertypes.StringValue `tfsdk:"id"`
	Name         supertypes.StringValue `tfsdk:"name"`
	Enabled      supertypes.BoolValue   `tfsdk:"enabled"`
	ProjectId    supertypes.StringValue `tfsdk:"project_id"`
}

func (m *ErrorMonitorDataSourceModel) Fill(ctx context.Context, data apiclient.ProjectMonitor) (diags diag.Diagnostics) {
	m.Id = supertypes.NewStringValue(data.Id)
	m.Name = supertypes.NewStringValue(data.Name)
	m.Enabled = supertypes.NewBoolValue(data.Enabled)
	m.ProjectId = supertypes.NewStringValue(data.ProjectId)
	return
}
