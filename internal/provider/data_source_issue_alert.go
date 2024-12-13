package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrytypes"
)

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
	nameStringAttribute := schema.StringAttribute{
		Computed: true,
	}
	intervalStringAttribute := schema.StringAttribute{
		MarkdownDescription: "Valid values are `1m`, `5m`, `15m`, `1h`, `1d`, `1w` and `30d` (`m` for minutes, `h` for hours, `d` for days, and `w` for weeks).",
		Computed:            true,
	}
	conditionComparisonTypeStringAttribute := schema.StringAttribute{
		MarkdownDescription: "Valid values are `count` and `percent`.",
		Computed:            true,
	}
	conditionComparisonIntervalStringAttribute := schema.StringAttribute{
		MarkdownDescription: "Valid values are `5m`, `15m`, `1h`, `1d`, `1w` and `30d` (`m` for minutes, `h` for hours, `d` for days, and `w` for weeks).",
		Computed:            true,
	}

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
			"conditions_v2": schema.ListNestedAttribute{
				MarkdownDescription: "A list of triggers that determine when the rule fires.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"first_seen_event": schema.SingleNestedAttribute{
							MarkdownDescription: "A new issue is created.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"regression_event": schema.SingleNestedAttribute{
							MarkdownDescription: "The issue changes state from resolved to unresolved.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"reappeared_event": schema.SingleNestedAttribute{
							MarkdownDescription: "The issue changes state from ignored to unresolved.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"new_high_priority_issue": schema.SingleNestedAttribute{
							MarkdownDescription: "Sentry marks a new issue as high priority.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"existing_high_priority_issue": schema.SingleNestedAttribute{
							MarkdownDescription: "Sentry marks an existing issue as high priority.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name": nameStringAttribute,
							},
						},
						"event_frequency": schema.SingleNestedAttribute{
							MarkdownDescription: "When the `comparison_type` is `count`, the number of events in an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the number of events in an issue is `value` % higher in `interval` compared to `comparison_interval` ago.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":                nameStringAttribute,
								"comparison_type":     conditionComparisonTypeStringAttribute,
								"comparison_interval": conditionComparisonIntervalStringAttribute,
								"value": schema.Int64Attribute{
									Computed: true,
								},
								"interval": intervalStringAttribute,
							},
						},
						"event_unique_user_frequency": schema.SingleNestedAttribute{
							MarkdownDescription: "When the `comparison_type` is `count`, the number of users affected by an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the number of users affected by an issue is `value` % higher in `interval` compared to `comparison_interval` ago.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":                nameStringAttribute,
								"comparison_type":     conditionComparisonTypeStringAttribute,
								"comparison_interval": conditionComparisonIntervalStringAttribute,
								"value": schema.Int64Attribute{
									Computed: true,
								},
								"interval": intervalStringAttribute,
							},
						},
						"event_frequency_percent": schema.SingleNestedAttribute{
							MarkdownDescription: "When the `comparison_type` is `count`, the percent of sessions affected by an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the percent of sessions affected by an issue is `value` % higher in `interval` compared to `comparison_interval` ago.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"name":                nameStringAttribute,
								"comparison_type":     conditionComparisonTypeStringAttribute,
								"comparison_interval": conditionComparisonIntervalStringAttribute,
								"value": schema.Float64Attribute{
									Computed: true,
								},
								"interval": schema.StringAttribute{
									MarkdownDescription: "Valid values are `5m`, `10m`, `30m`, and `1h` (`m` for minutes, `h` for hours).",
									Computed:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("5m", "10m", "30m", "1h"),
									},
								},
							},
						},
					},
				},
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
	var data IssueAlertModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := d.apiClient.GetProjectRuleWithResponse(
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
