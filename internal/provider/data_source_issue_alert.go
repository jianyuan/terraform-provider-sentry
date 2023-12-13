package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ datasource.DataSource = &IssueAlertDataSource{}

func NewIssueAlertDataSource() datasource.DataSource {
	return &IssueAlertDataSource{}
}

type IssueAlertDataSource struct {
	client *sentry.Client
}

type IssueAlertDataSourceModel struct {
	Id           types.String         `tfsdk:"id"`
	Organization types.String         `tfsdk:"organization"`
	Project      types.String         `tfsdk:"project"`
	Name         types.String         `tfsdk:"name"`
	Conditions   jsontypes.Normalized `tfsdk:"conditions"`
	Filters      jsontypes.Normalized `tfsdk:"filters"`
	Actions      jsontypes.Normalized `tfsdk:"actions"`
	ActionMatch  types.String         `tfsdk:"action_match"`
	FilterMatch  types.String         `tfsdk:"filter_match"`
	Frequency    types.Int64          `tfsdk:"frequency"`
	Environment  types.String         `tfsdk:"environment"`
	Owner        types.String         `tfsdk:"owner"`
}

func (m *IssueAlertDataSourceModel) Fill(organization string, alert sentry.IssueAlert) error {
	m.Id = types.StringPointerValue(alert.ID)
	m.Organization = types.StringValue(organization)
	m.Project = types.StringValue(alert.Projects[0])
	m.Name = types.StringPointerValue(alert.Name)
	m.ActionMatch = types.StringPointerValue(alert.ActionMatch)
	m.FilterMatch = types.StringPointerValue(alert.FilterMatch)
	m.Owner = types.StringPointerValue(alert.Owner)

	// Remove the name from the conditions, filters, and actions. They are added by the API.
	// We do this to avoid a diff when the user updates the resource.
	for _, m := range alert.Conditions {
		delete(m, "name")
	}
	for _, m := range alert.Filters {
		delete(m, "name")
	}
	for _, m := range alert.Actions {
		delete(m, "name")
	}

	if conditions, err := json.Marshal(alert.Conditions); err == nil {
		m.Conditions = jsontypes.NewNormalizedValue(string(conditions))
	} else {
		m.Conditions = jsontypes.NewNormalizedNull()
	}

	if filters, err := json.Marshal(alert.Filters); err == nil {
		m.Filters = jsontypes.NewNormalizedValue(string(filters))
	} else {
		m.Filters = jsontypes.NewNormalizedNull()
	}

	if actions, err := json.Marshal(alert.Actions); err == nil {
		m.Actions = jsontypes.NewNormalizedValue(string(actions))
	} else {
		m.Actions = jsontypes.NewNormalizedNull()
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
			"organization": schema.StringAttribute{
				MarkdownDescription: "The slug of the organization the resource belongs to.",
				Required:            true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The slug of the project the resource belongs to.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The issue alert name.",
				Computed:            true,
			},
			"conditions": schema.StringAttribute{
				MarkdownDescription: "List of conditions.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"filters": schema.StringAttribute{
				MarkdownDescription: "A list of filters that determine if a rule fires after the necessary conditions have been met.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"actions": schema.StringAttribute{
				MarkdownDescription: "List of actions.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
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

func (d *IssueAlertDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Issue alert not found: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading issue alert: %s", err.Error()))
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *action); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error filling issue alert: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
