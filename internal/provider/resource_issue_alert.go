package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

var _ resource.Resource = &IssueAlertResource{}
var _ resource.ResourceWithImportState = &IssueAlertResource{}
var _ resource.ResourceWithUpgradeState = &IssueAlertResource{}

func NewIssueAlertResource() resource.Resource {
	return &IssueAlertResource{}
}

type IssueAlertResource struct {
	client *sentry.Client
}

type IssueAlertResourceModel struct {
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

func (m *IssueAlertResourceModel) Fill(organization string, alert sentry.IssueAlert) error {
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

func (r *IssueAlertResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_issue_alert"
}

func (r *IssueAlertResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create an Issue Alert Rule for a Project. See the [Sentry Documentation](https://docs.sentry.io/api/alerts/create-an-issue-alert-rule-for-a-project/) for more information.",

		Version: 2,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"conditions": schema.StringAttribute{
				MarkdownDescription: "List of conditions.",
				Required:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"filters": schema.StringAttribute{
				MarkdownDescription: "A list of filters that determine if a rule fires after the necessary conditions have been met.",
				Optional:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"actions": schema.StringAttribute{
				MarkdownDescription: "List of actions.",
				Required:            true,
				CustomType:          jsontypes.NormalizedType{},
			},
			"action_match": schema.StringAttribute{
				MarkdownDescription: "Trigger actions when an event is captured by Sentry and `any` or `all` of the specified conditions happen.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("all", "any"),
				},
			},
			"filter_match": schema.StringAttribute{
				MarkdownDescription: "A string determining which filters need to be true before any actions take place. Required when a value is provided for `filters`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("all", "any", "none"),
				},
			},
			"frequency": schema.Int64Attribute{
				MarkdownDescription: "Perform actions at most once every `X` minutes for this issue.",
				Required:            true,
			},
			"environment": schema.StringAttribute{
				MarkdownDescription: "Perform issue alert in a specific environment.",
				Optional:            true,
			},
			"owner": schema.StringAttribute{
				MarkdownDescription: "The ID of the team or user that owns the rule.",
				Optional:            true,
			},
		},
	}
}

func (r *IssueAlertResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sentry.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sentry.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *IssueAlertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IssueAlertResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &sentry.IssueAlert{
		Name:        data.Name.ValueStringPointer(),
		ActionMatch: data.ActionMatch.ValueStringPointer(),
		FilterMatch: data.FilterMatch.ValueStringPointer(),
		Frequency:   sentry.JsonNumber(json.Number(data.Frequency.String())),
		Owner:       data.Owner.ValueStringPointer(),
		Environment: data.Environment.ValueStringPointer(),
		Projects:    []string{data.Project.String()},
	}
	resp.Diagnostics.Append(data.Conditions.Unmarshal(&params.Conditions)...)
	resp.Diagnostics.Append(data.Filters.Unmarshal(&params.Filters)...)
	resp.Diagnostics.Append(data.Actions.Unmarshal(&params.Actions)...)

	if resp.Diagnostics.HasError() {
		return
	}

	action, _, err := r.client.IssueAlerts.Create(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		params,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating issue alert: %s", err.Error()))
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *action); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error filling issue alert: %s", err.Error()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IssueAlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IssueAlertResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	action, apiResp, err := r.client.IssueAlerts.Get(
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

func (r *IssueAlertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IssueAlertResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := &sentry.IssueAlert{
		Name:        data.Name.ValueStringPointer(),
		ActionMatch: data.ActionMatch.ValueStringPointer(),
		FilterMatch: data.FilterMatch.ValueStringPointer(),
		Frequency:   sentry.JsonNumber(json.Number(data.Frequency.String())),
		Owner:       data.Owner.ValueStringPointer(),
		Environment: data.Environment.ValueStringPointer(),
		Projects:    []string{data.Project.String()},
	}
	resp.Diagnostics.Append(data.Conditions.Unmarshal(&params.Conditions)...)
	resp.Diagnostics.Append(data.Filters.Unmarshal(&params.Filters)...)
	resp.Diagnostics.Append(data.Actions.Unmarshal(&params.Actions)...)

	if resp.Diagnostics.HasError() {
		return
	}

	action, apiResp, err := r.client.IssueAlerts.Update(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
		params,
	)
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Notification Action not found: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error updating notification action: %s", err.Error()))
		return
	}

	if err := data.Fill(data.Organization.ValueString(), *action); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error filling issue alert: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IssueAlertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IssueAlertResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.IssueAlerts.Delete(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error deleting issue alert: %s", err.Error()))
		return
	}
}

func (r *IssueAlertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, project, actionId, err := splitThreePartID(req.ID, "organization", "project-slug", "alert-id")
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), organization,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("project"), project,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), actionId,
	)...)
}

func (r *IssueAlertResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var data IssueAlertResourceModel

				resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

				if resp.Diagnostics.HasError() {
					return
				}

				data.Id = types.StringValue(buildThreePartID(data.Organization.ValueString(), data.Project.ValueString(), data.Id.ValueString()))

				resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			},
		},
		1: {
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var data IssueAlertResourceModel

				resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

				if resp.Diagnostics.HasError() {
					return
				}

				organization, project, alertId, err := splitThreePartID(data.Id.ValueString(), "organization", "project-slug", "alert-id")
				if err != nil {
					resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
					return
				}

				data.Id = types.StringValue(alertId)
				data.Organization = types.StringValue(organization)
				data.Project = types.StringValue(project)

				resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			},
		},
	}
}
