package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
)

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

func (m *NotificationActionResourceModel) Fill(action sentry.NotificationAction, projectIdToSlugMap map[string]string) error {
	if action.ID == nil {
		m.Id = types.StringNull()
	} else {
		m.Id = types.StringValue(action.ID.String())
	}

	m.TriggerType = types.StringPointerValue(action.TriggerType)
	m.ServiceType = types.StringPointerValue(action.ServiceType)

	if action.IntegrationId == nil {
		m.IntegrationId = types.StringNull()
	} else {
		m.IntegrationId = types.StringValue(action.IntegrationId.String())
	}

	switch targetIdentifier := action.TargetIdentifier.(type) {
	case string:
		m.TargetIdentifier = types.StringValue(targetIdentifier)
	case int64:
		m.TargetIdentifier = types.StringValue(strconv.FormatInt(targetIdentifier, 10))
	case nil:
		m.TargetIdentifier = types.StringNull()
	}

	m.TargetDisplay = types.StringPointerValue(action.TargetDisplay)

	if len(action.Projects) > 0 {
		projectElements := []attr.Value{}
		for _, projectId := range action.Projects {
			projectElements = append(projectElements, types.StringValue(projectIdToSlugMap[projectId.String()]))
		}

		m.Projects = types.ListValueMust(types.StringType, projectElements)
	}

	return nil
}

var _ resource.Resource = &NotificationActionResource{}
var _ resource.ResourceWithConfigure = &NotificationActionResource{}
var _ resource.ResourceWithImportState = &NotificationActionResource{}

func NewNotificationActionResource() resource.Resource {
	return &NotificationActionResource{}
}

type NotificationActionResource struct {
	baseResource
}

func (r *NotificationActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_action"
}

func (r *NotificationActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create a Spike Protection Notification Action. See the [Sentry Documentation](https://docs.sentry.io/api/alerts/create-a-spike-protection-notification-action/) for more information.",

		Attributes: map[string]schema.Attribute{
			"id":           ResourceIdAttribute(),
			"organization": ResourceOrganizationAttribute(),
			"trigger_type": schema.StringAttribute{
				Description: "The type of trigger that will activate this action. Valid values are `spike-protection`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("spike-protection"),
				},
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

func (r *NotificationActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NotificationActionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projects := []string{}
	resp.Diagnostics.Append(data.Projects.ElementsAs(ctx, &projects, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
		diagutils.AddClientError(resp.Diagnostics, "create", err)
		return
	}

	var projectIdToSlugMap map[string]string
	if len(action.Projects) > 0 {
		projectIdToSlugMap, err = sentryclient.GetProjectIdToSlugMap(ctx, r.client)
		if err != nil {
			diagutils.AddClientError(resp.Diagnostics, "read projects", err)
			return
		}
	}

	if err := data.Fill(*action, projectIdToSlugMap); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
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
		diagutils.AddNotFoundError(resp.Diagnostics, "notification action")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "read", err)
		return
	}

	var projectIdToSlugMap map[string]string
	if len(action.Projects) > 0 {
		projectIdToSlugMap, err = sentryclient.GetProjectIdToSlugMap(ctx, r.client)
		if err != nil {
			diagutils.AddClientError(resp.Diagnostics, "read projects", err)
			return
		}
	}

	if err := data.Fill(*action, projectIdToSlugMap); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
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
	if resp.Diagnostics.HasError() {
		return
	}

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
		diagutils.AddNotFoundError(resp.Diagnostics, "notification action")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "update", err)
		return
	}

	var projectIdToSlugMap map[string]string
	if len(action.Projects) > 0 {
		projectIdToSlugMap, err = sentryclient.GetProjectIdToSlugMap(ctx, r.client)
		if err != nil {
			diagutils.AddClientError(resp.Diagnostics, "read projects", err)
			return
		}
	}

	if err := data.Fill(*action, projectIdToSlugMap); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
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
		diagutils.AddClientError(resp.Diagnostics, "delete", err)
		return
	}
}

func (r *NotificationActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	org, actionId, err := splitTwoPartID(req.ID, "organization", "action-id")
	if err != nil {
		diagutils.AddImportError(resp.Diagnostics, err)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), org,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), actionId,
	)...)
}
