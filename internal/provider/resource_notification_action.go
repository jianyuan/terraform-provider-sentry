package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
)

var _ resource.Resource = &NotificationActionResource{}
var _ resource.ResourceWithImportState = &NotificationActionResource{}

func NewNotificationActionResource() resource.Resource {
	return &NotificationActionResource{}
}

type NotificationActionResource struct {
	client *sentry.Client
}

type NotificationActionResourceModel struct {
	Id               types.String `tfsdk:"id"`
	Organization     types.String `tfsdk:"organization"`
	TriggerType      types.String `tfsdk:"trigger_type"`
	ServiceType      types.String `tfsdk:"service_type"`
	IntegrationId    types.String `tfsdk:"integration_id"`
	TargetIdentifier types.String `tfsdk:"target_identifier"`
	TargetDisplay    types.String `tfsdk:"target_display"`
	Projects         types.List   `tfsdk:"projects"`
}

func (r *NotificationActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_action"
}

func (r *NotificationActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a Spike Protection Notification Action. See the [Sentry Documentation](https://docs.sentry.io/api/alerts/create-a-spike-protection-notification-action/) for more information.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of this resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization": schema.StringAttribute{
				Description: "The slug of the organization the project belongs to.",
				Required:    true,
			},
			"trigger_type": schema.StringAttribute{
				Description: "The type of trigger that will activate this action. Valid values are `spike_protection`.",
				Required:    true,
			},
			"service_type": schema.StringAttribute{
				Description: "The service that is used for sending the notification.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("email", "slack", "sentry_notification", "pagerduty", "opsgenie"),
				},
			},
			"integration_id": schema.StringAttribute{
				Description: "The ID of the integration that is used for sending the notification. Use the `sentry_organization_integration` data source to retrieve an integration. Required if `service_type` is `slack`, `pagerduty` or `opsgenie`.",
				Optional:    true,
			},
			"target_identifier": schema.StringAttribute{
				Description: "The identifier of the target that is used for sending the notification (e.g. Slack channel ID). Required if `service_type` is `slack` or `opsgenie`.",
				Optional:    true,
			},
			"target_display": schema.StringAttribute{
				Description: "The display name of the target that is used for sending the notification (e.g. Slack channel name). Required if `service_type` is `slack` or `opsgenie`.",
				Optional:    true,
			},
			"projects": schema.ListAttribute{
				Description: "The list of project slugs that the Notification Action is created for.",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *NotificationActionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NotificationActionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projects := []string{}
	resp.Diagnostics.Append(data.Projects.ElementsAs(ctx, &projects, false)...)
	action, _, err := r.client.NotificationActions.Create(
		ctx,
		data.Organization.ValueString(),
		&sentry.CreateNotificationActionParams{
			TriggerType:      data.TriggerType.ValueStringPointer(),
			ServiceType:      data.ServiceType.ValueStringPointer(),
			IntegrationId:    (*json.Number)(data.IntegrationId.ValueStringPointer()),
			TargetIdentifier: data.TargetIdentifier.ValueStringPointer(),
			TargetDisplay:    data.TargetDisplay.ValueStringPointer(),
			Projects:         projects,
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating notification action: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(r.mapToModel(ctx, action, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotificationActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NotificationActionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	action, apiResp, err := r.client.NotificationActions.Get(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Notification Action not found: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error reading notification action: %s", err.Error()))
		return
	}

	resp.Diagnostics.Append(r.mapToModel(ctx, action, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotificationActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NotificationActionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projects := []string{}
	resp.Diagnostics.Append(data.Projects.ElementsAs(ctx, &projects, false)...)
	action, apiResp, err := r.client.NotificationActions.Update(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
		&sentry.UpdateNotificationActionParams{
			TriggerType:      data.TriggerType.ValueStringPointer(),
			ServiceType:      data.ServiceType.ValueStringPointer(),
			IntegrationId:    (*json.Number)(data.IntegrationId.ValueStringPointer()),
			TargetIdentifier: data.TargetIdentifier.ValueStringPointer(),
			TargetDisplay:    data.TargetDisplay.ValueStringPointer(),
			Projects:         projects,
		},
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

	resp.Diagnostics.Append(r.mapToModel(ctx, action, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotificationActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NotificationActionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := r.client.NotificationActions.Delete(
		ctx,
		data.Organization.ValueString(),
		data.Id.ValueString(),
	)
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error disabling spike protection: %s", err.Error()))
		return
	}
}

func (r *NotificationActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	org, actionId, err := splitTwoPartID(req.ID, "organization", "action-id")
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), org,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), actionId,
	)...)
}

func (r *NotificationActionResource) mapToModel(ctx context.Context, action *sentry.NotificationAction, data *NotificationActionResourceModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	data.Id = types.StringPointerValue((*string)(action.ID))
	data.TriggerType = types.StringPointerValue(action.TriggerType)
	data.ServiceType = types.StringPointerValue(action.ServiceType)
	data.IntegrationId = types.StringPointerValue((*string)(action.IntegrationId))
	switch targetIdentifier := action.TargetIdentifier.(type) {
	case string:
		data.TargetIdentifier = types.StringValue(targetIdentifier)
	case int64:
		data.TargetIdentifier = types.StringValue(strconv.FormatInt(targetIdentifier, 10))
	case nil:
		data.TargetIdentifier = types.StringNull()
	}
	data.TargetDisplay = types.StringPointerValue(action.TargetDisplay)

	if len(action.Projects) > 0 {
		projectIdToSlugMap, err := sentryclient.GetProjectIdToSlugMap(ctx, r.client)
		if err != nil {
			diagnostics.AddError("Client Error", fmt.Sprintf("Error reading projects: %s", err.Error()))
			return diagnostics
		}

		projectElements := []attr.Value{}
		for _, projectId := range action.Projects {
			projectElements = append(projectElements, types.StringValue(projectIdToSlugMap[projectId.String()]))
		}

		projects, diags := types.ListValue(types.StringType, projectElements)
		data.Projects = projects
		diagnostics.Append(diags...)
	}

	return diagnostics
}
