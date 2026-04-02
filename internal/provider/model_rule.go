package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type SourceSentryRuleModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Project      types.String `tfsdk:"project"`
	Name         types.String `tfsdk:"name"`
	ActionMatch  types.String `tfsdk:"action_match"`
	FilterMatch  types.String `tfsdk:"filter_match"`
	Actions      types.List   `tfsdk:"actions"`
	Conditions   types.List   `tfsdk:"conditions"`
	Filters      types.List   `tfsdk:"filters"`
	Frequency    types.Int64  `tfsdk:"frequency"`
	Environment  types.String `tfsdk:"environment"`
}
