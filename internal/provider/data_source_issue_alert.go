package provider

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
)

type IssueAlertDataSourceModel struct {
	Id           types.String          `tfsdk:"id"`
	Organization types.String          `tfsdk:"organization"`
	Project      types.String          `tfsdk:"project"`
	Name         types.String          `tfsdk:"name"`
	Conditions   sentrytypes.LossyJson `tfsdk:"conditions"`
	Filters      sentrytypes.LossyJson `tfsdk:"filters"`
	Actions      sentrytypes.LossyJson `tfsdk:"actions"`
	ActionMatch  types.String          `tfsdk:"action_match"`
	FilterMatch  types.String          `tfsdk:"filter_match"`
	Frequency    types.Int64           `tfsdk:"frequency"`
	Environment  types.String          `tfsdk:"environment"`
	Owner        types.String          `tfsdk:"owner"`
}

func (m *IssueAlertDataSourceModel) Fill(organization string, alert sentry.IssueAlert) error {
	m.Id = types.StringPointerValue(alert.ID)
	m.Organization = types.StringValue(organization)
	m.Project = types.StringValue(alert.Projects[0])
	m.Name = types.StringPointerValue(alert.Name)
	m.ActionMatch = types.StringPointerValue(alert.ActionMatch)
	m.FilterMatch = types.StringPointerValue(alert.FilterMatch)
	m.Owner = types.StringPointerValue(alert.Owner)

	m.Conditions = sentrytypes.NewLossyJsonNull()
	if len(alert.Conditions) > 0 {
		if conditions, err := json.Marshal(alert.Conditions); err == nil {
			m.Conditions = sentrytypes.NewLossyJsonValue(string(conditions))
		} else {
			return err
		}
	}

	m.Filters = sentrytypes.NewLossyJsonNull()
	if len(alert.Filters) > 0 {
		if filters, err := json.Marshal(alert.Filters); err == nil {
			m.Filters = sentrytypes.NewLossyJsonValue(string(filters))
		} else {
			return err
		}
	}

	m.Actions = sentrytypes.NewLossyJsonNull()
	if len(alert.Actions) > 0 {
		if actions, err := json.Marshal(alert.Actions); err == nil && len(actions) > 0 {
			m.Actions = sentrytypes.NewLossyJsonValue(string(actions))
		} else {
			return err
		}
	}

	frequency, err := alert.Frequency.Int64()
	if err != nil {
		return err
	}
	m.Frequency = types.Int64Value(frequency)

	m.Environment = types.StringPointerValue(alert.Environment)
	m.Owner = types.StringPointerValue(alert.Owner)

	return nil
}

var _ datasource.DataSource = &IssueAlertDataSource{}
var _ datasource.DataSourceWithConfigure = &IssueAlertDataSource{}

func NewIssueAlertDataSource() datasource.DataSource {
	return &IssueAlertDataSource{}
}

type IssueAlertDataSource struct {
	baseDataSource
}

func (d *IssueAlertDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_issue_alert"
}

func (d *IssueAlertDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Sentry Issue Alert data source. See the [Sentry documentation](https://docs.sentry.io/api/alerts/retrieve-an-issue-alert-rule-for-a-project/) for more information.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Required:    true,
			},
			"organization": DataSourceOrganizationAttribute(),
			"project":      DataSourceProjectAttribute(),
			"name": schema.StringAttribute{
				MarkdownDescription: "The issue alert name.",
				Computed:            true,
			},
			"conditions": schema.StringAttribute{
				MarkdownDescription: "List of conditions. In JSON string format.",
				Computed:            true,
				CustomType:          sentrytypes.LossyJsonType{},
			},
			"filters": schema.StringAttribute{
				MarkdownDescription: "A list of filters that determine if a rule fires after the necessary conditions have been met. In JSON string format.",
				Computed:            true,
				CustomType:          sentrytypes.LossyJsonType{},
			},
			"actions": schema.StringAttribute{
				MarkdownDescription: "List of actions. In JSON string format.",
				Computed:            true,
				CustomType:          sentrytypes.LossyJsonType{},
			},
			"action_match": schema.StringAttribute{
				MarkdownDescription: "Trigger actions when an event is captured by Sentry and `any` or `all` of the specified conditions happen.",
				Computed:            true,
			},
			"filter_match": schema.StringAttribute{
				MarkdownDescription: "A string determining which filters need to be true before any actions take place. Required when a value is provided for `filters`.",
				Computed:            true,
			},
			"frequency": schema.Int64Attribute{
				MarkdownDescription: "Perform actions at most once every `X` minutes for this issue.",
				Computed:            true,
			},
			"environment": schema.StringAttribute{
				MarkdownDescription: "Perform issue alert in a specific environment.",
				Computed:            true,
			},
			"owner": schema.StringAttribute{
				MarkdownDescription: "The ID of the team or user that owns the rule.",
				Computed:            true,
			},
		},
	}
}

func (d *IssueAlertDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IssueAlertResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	action, apiResp, err := d.client.IssueAlerts.Get(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		diagutils.AddNotFoundError(resp.Diagnostics, "issue alert")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "read", err)
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *action); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
