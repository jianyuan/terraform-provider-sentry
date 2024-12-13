package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/go-utils/sliceutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

var _ resource.Resource = &IssueAlertResource{}
var _ resource.ResourceWithConfigure = &IssueAlertResource{}
var _ resource.ResourceWithImportState = &IssueAlertResource{}
var _ resource.ResourceWithUpgradeState = &IssueAlertResource{}

func NewIssueAlertResource() resource.Resource {
	return &IssueAlertResource{}
}

type IssueAlertResource struct {
	baseResource
}

func (r *IssueAlertResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_issue_alert"
}

func (r *IssueAlertResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// idStringAttribute := schema.StringAttribute{
	// 	Computed: true,
	// }
	nameStringAttribute := schema.StringAttribute{
		Computed: true,
	}
	intervalStringAttribute := schema.StringAttribute{
		MarkdownDescription: "Valid values are `1m`, `5m`, `15m`, `1h`, `1d`, `1w` and `30d` (`m` for minutes, `h` for hours, `d` for days, and `w` for weeks).",
		Optional:            true,
		Validators: []validator.String{
			stringvalidator.OneOf("1m", "5m", "15m", "1h", "1d", "1w", "30d"),
		},
	}
	conditionComparisonTypeStringAttribute := schema.StringAttribute{
		MarkdownDescription: "Valid values are `count` and `percent`.",
		Required:            true,
		Validators: []validator.String{
			stringvalidator.OneOf("count", "percent"),
		},
	}
	conditionComparisonIntervalStringAttribute := schema.StringAttribute{
		MarkdownDescription: "Valid values are `5m`, `15m`, `1h`, `1d`, `1w` and `30d` (`m` for minutes, `h` for hours, `d` for days, and `w` for weeks).",
		Optional:            true,
		Validators: []validator.String{
			stringvalidator.OneOf("5m", "15m", "1h", "1d", "1w", "30d"),
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: `Create an Issue Alert Rule for a Project. See the [Sentry Documentation](https://docs.sentry.io/api/alerts/create-an-issue-alert-rule-for-a-project/) for more information.

TODO

Please note the following changes since v0.12.0:
- The attributes ` + "`conditions`" + `, ` + "`filters`" + `, and ` + "`actions`" + ` are in JSON string format. The types must match the Sentry API, otherwise Terraform will incorrectly detect a drift. Use ` + "`parseint(\"string\", 10)`" + ` to convert a string to an integer. Avoid using ` + "`jsonencode()`" + ` as it is unable to distinguish between an integer and a float.
- The attribute ` + "`internal_id`" + ` has been removed. Use ` + "`id`" + ` instead.
- The attribute ` + "`id`" + ` is now the ID of the issue alert. Previously, it was a combination of the organization, project, and issue alert ID.
		`,

		Version: 2,

		Attributes: map[string]schema.Attribute{
			"id":           ResourceIdAttribute(),
			"organization": ResourceOrganizationAttribute(),
			"project":      ResourceProjectAttribute(),
			"name": schema.StringAttribute{
				MarkdownDescription: "The issue alert name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"conditions": schema.StringAttribute{
				MarkdownDescription: "**Deprecated** in favor of `condition`. A list of triggers that determine when the rule fires. In JSON string format.",
				DeprecationMessage:  "Use `condition` instead.",
				Optional:            true,
				CustomType: sentrytypes.LossyJsonType{
					IgnoreKeys: []string{"name"},
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("conditions_v2")),
				},
			},
			"conditions_v2": schema.ListNestedAttribute{
				MarkdownDescription: "A list of triggers that determine when the rule fires.",
				Optional:            true,
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("conditions")),
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: tfutils.WithMutuallyExclusiveValidator(map[string]schema.SingleNestedAttribute{
						"first_seen_event": {
							MarkdownDescription: "A new issue is created.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"regression_event": {
							MarkdownDescription: "The issue changes state from resolved to unresolved.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"reappeared_event": {
							MarkdownDescription: "The issue changes state from ignored to unresolved.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"new_high_priority_issue": {
							MarkdownDescription: "Sentry marks a new issue as high priority.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"existing_high_priority_issue": {
							MarkdownDescription: "Sentry marks an existing issue as high priority.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"event_frequency": {
							MarkdownDescription: "When the `comparison_type` is `count`, the number of events in an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the number of events in an issue is `value` % higher in `interval` compared to `comparison_interval` ago.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name":                nameStringAttribute,
								"comparison_type":     conditionComparisonTypeStringAttribute,
								"comparison_interval": conditionComparisonIntervalStringAttribute,
								"value": schema.Int64Attribute{
									Required: true,
								},
								"interval": intervalStringAttribute,
							},
						},
						"event_unique_user_frequency": {
							MarkdownDescription: "When the `comparison_type` is `count`, the number of users affected by an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the number of users affected by an issue is `value` % higher in `interval` compared to `comparison_interval` ago.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name":                nameStringAttribute,
								"comparison_type":     conditionComparisonTypeStringAttribute,
								"comparison_interval": conditionComparisonIntervalStringAttribute,
								"value": schema.Int64Attribute{
									Required: true,
								},
								"interval": intervalStringAttribute,
							},
						},
						"event_frequency_percent": {
							MarkdownDescription: "When the `comparison_type` is `count`, the percent of sessions affected by an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the percent of sessions affected by an issue is `value` % higher in `interval` compared to `comparison_interval` ago.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"name":                nameStringAttribute,
								"comparison_type":     conditionComparisonTypeStringAttribute,
								"comparison_interval": conditionComparisonIntervalStringAttribute,
								"value": schema.Float64Attribute{
									Required: true,
								},
								"interval": schema.StringAttribute{
									MarkdownDescription: "Valid values are `5m`, `10m`, `30m`, and `1h` (`m` for minutes, `h` for hours).",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("5m", "10m", "30m", "1h"),
									},
								},
							},
						},
					}),
				},
			},
			"filters": schema.StringAttribute{
				MarkdownDescription: "A list of filters that determine if a rule fires after the necessary conditions have been met. In JSON string format.",
				Optional:            true,
				CustomType:          sentrytypes.LossyJsonType{},
			},
			"actions": schema.StringAttribute{
				MarkdownDescription: "List of actions. In JSON string format.",
				Required:            true,
				CustomType:          sentrytypes.LossyJsonType{},
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

func (r *IssueAlertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IssueAlertModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.CreateProjectRuleJSONRequestBody{
		Name:        data.Name.ValueString(),
		ActionMatch: data.ActionMatch.ValueString(),
		FilterMatch: data.FilterMatch.ValueString(),
		Frequency:   data.Frequency.ValueInt64(),
		Owner:       data.Owner.ValueStringPointer(),
		Environment: data.Environment.ValueStringPointer(),
		Projects:    []string{data.Project.ValueString()},
	}

	if !data.Conditions.IsNull() {
		resp.Diagnostics.Append(data.Conditions.Unmarshal(&body.Conditions)...)
	} else if data.ConditionsV2 != nil {
		body.Conditions = sliceutils.Map(func(item IssueAlertConditionModel) apiclient.ProjectRuleCondition {
			return item.ToApi()
		}, *data.ConditionsV2)
	} else {
		panic("provider error: conditions is required")
	}

	if !data.Filters.IsNull() {
		resp.Diagnostics.Append(data.Filters.Unmarshal(&body.Filters)...)
	} else {
		body.Filters = []map[string]interface{}{}
	}

	if !data.Actions.IsNull() {
		resp.Diagnostics.Append(data.Actions.Unmarshal(&body.Actions)...)
	} else {
		body.Actions = []map[string]interface{}{}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CreateProjectRuleWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		body,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("create", err))
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("create", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IssueAlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IssueAlertModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GetProjectRuleWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("issue alert"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IssueAlertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IssueAlertModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := apiclient.UpdateProjectRuleJSONRequestBody{
		Name:        data.Name.ValueString(),
		ActionMatch: data.ActionMatch.ValueString(),
		FilterMatch: data.FilterMatch.ValueString(),
		Frequency:   data.Frequency.ValueInt64(),
		Owner:       data.Owner.ValueStringPointer(),
		Environment: data.Environment.ValueStringPointer(),
		Projects:    []string{data.Project.ValueString()},
	}

	if !data.Conditions.IsNull() {
		resp.Diagnostics.Append(data.Conditions.Unmarshal(&body.Conditions)...)
	} else if data.ConditionsV2 != nil {
		body.Conditions = sliceutils.Map(func(item IssueAlertConditionModel) apiclient.ProjectRuleCondition {
			return item.ToApi()
		}, *data.ConditionsV2)
	} else {
		panic("provider error: conditions is required")
	}

	if !data.Filters.IsNull() {
		resp.Diagnostics.Append(data.Filters.Unmarshal(&body.Filters)...)
	} else {
		body.Filters = []map[string]interface{}{}
	}

	if !data.Actions.IsNull() {
		resp.Diagnostics.Append(data.Actions.Unmarshal(&body.Actions)...)
	} else {
		body.Actions = []map[string]interface{}{}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.UpdateProjectRuleWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
		body,
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("issue alert"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("update", httpResp.StatusCode(), httpResp.Body))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *httpResp.JSON200)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IssueAlertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IssueAlertModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DeleteProjectRuleWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.Project.ValueString(),
		data.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		return
	} else if httpResp.StatusCode() != http.StatusAccepted {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("delete", httpResp.StatusCode(), httpResp.Body))
		return
	}
}

func (r *IssueAlertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tfutils.ImportStateThreePartId(ctx, "organization", "project", req, resp)
}

func (r *IssueAlertResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	type modelV0 struct {
		Id           types.String `tfsdk:"id"`
		Organization types.String `tfsdk:"organization"`
		Project      types.String `tfsdk:"project"`
		Name         types.String `tfsdk:"name"`
		Conditions   types.List   `tfsdk:"conditions"`
		Filters      types.List   `tfsdk:"filters"`
		Actions      types.List   `tfsdk:"actions"`
		ActionMatch  types.String `tfsdk:"action_match"`
		FilterMatch  types.String `tfsdk:"filter_match"`
		Frequency    types.Int64  `tfsdk:"frequency"`
		Environment  types.String `tfsdk:"environment"`
	}

	return map[int64]resource.StateUpgrader{
		0: {
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				// No-op
			},
		},
		1: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
					},
					"organization": schema.StringAttribute{
						Required: true,
					},
					"project": schema.StringAttribute{
						Required: true,
					},
					"name": schema.StringAttribute{
						Required: true,
					},
					"conditions": schema.ListAttribute{
						ElementType: types.MapType{
							ElemType: types.StringType,
						},
						Required: true,
					},
					"filters": schema.ListAttribute{
						ElementType: types.MapType{
							ElemType: types.StringType,
						},
						Optional: true,
					},
					"actions": schema.ListAttribute{
						ElementType: types.MapType{
							ElemType: types.StringType,
						},
						Required: true,
					},
					"action_match": schema.StringAttribute{
						Optional: true,
					},
					"filter_match": schema.StringAttribute{
						Optional: true,
					},
					"frequency": schema.Int64Attribute{
						Optional: true,
					},
					"environment": schema.StringAttribute{
						Optional: true,
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var priorStateData modelV0

				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)
				if resp.Diagnostics.HasError() {
					return
				}

				organization, project, actionId, err := splitThreePartID(priorStateData.Id.ValueString(), "organization", "project-slug", "alert-id")
				if err != nil {
					resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Error parsing ID: %s", err.Error()))
					return
				}

				upgradedStateData := IssueAlertModel{
					Id:           types.StringValue(actionId),
					Organization: types.StringValue(organization),
					Project:      types.StringValue(project),
					Name:         priorStateData.Name,
					ActionMatch:  priorStateData.ActionMatch,
					FilterMatch:  priorStateData.FilterMatch,
					Frequency:    priorStateData.Frequency,
					Environment:  priorStateData.Environment,
				}

				upgradedStateData.Conditions = sentrytypes.NewLossyJsonNull()
				if !priorStateData.Conditions.IsNull() {
					conditions := []map[string]string{}
					resp.Diagnostics.Append(priorStateData.Conditions.ElementsAs(ctx, &conditions, false)...)
					if resp.Diagnostics.HasError() {
						return
					}

					if len(conditions) > 0 {
						upgradedStateData.Conditions = sentrytypes.NewLossyJsonValue(string(must.Get(json.Marshal(conditions))))
					}
				}

				upgradedStateData.Filters = sentrytypes.NewLossyJsonNull()
				if !priorStateData.Filters.IsNull() {
					filters := []map[string]string{}
					resp.Diagnostics.Append(priorStateData.Filters.ElementsAs(ctx, &filters, false)...)
					if resp.Diagnostics.HasError() {
						return
					}

					if len(filters) > 0 {
						upgradedStateData.Filters = sentrytypes.NewLossyJsonValue(string(must.Get(json.Marshal(filters))))
					}
				}

				upgradedStateData.Actions = sentrytypes.NewLossyJsonNull()
				if !priorStateData.Actions.IsNull() {
					actions := []map[string]string{}
					resp.Diagnostics.Append(priorStateData.Actions.ElementsAs(ctx, &actions, false)...)
					if resp.Diagnostics.HasError() {
						return
					}

					if len(actions) > 0 {
						upgradedStateData.Actions = sentrytypes.NewLossyJsonValue(string(must.Get(json.Marshal(actions))))
					}
				}

				resp.Diagnostics.Append(resp.State.Set(ctx, &upgradedStateData)...)
			},
		},
	}
}
