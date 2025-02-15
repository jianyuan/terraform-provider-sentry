package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
)

type MonitorConfigScheduleIntervalResourceModel struct {
	Year   types.Int64 `tfsdk:"year"`
	Month  types.Int64 `tfsdk:"month"`
	Week   types.Int64 `tfsdk:"week"`
	Day    types.Int64 `tfsdk:"day"`
	Hour   types.Int64 `tfsdk:"hour"`
	Minute types.Int64 `tfsdk:"minute"`
}

func (m MonitorConfigScheduleIntervalResourceModel) SchemaAttributes() map[string]schema.Attribute {
	attributeNames := []string{"year", "month", "week", "day", "hour", "minute"}

	attributes := make(map[string]schema.Attribute, len(attributeNames))

	for _, name := range attributeNames {
		var conflictingPaths []path.Expression

		for _, conflictingName := range attributeNames {
			if conflictingName != name {
				conflictingPaths = append(conflictingPaths, path.MatchRelative().AtParent().AtName(conflictingName))
			}
		}

		attributes[name] = schema.Int64Attribute{
			Optional: true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
				int64validator.ConflictsWith(conflictingPaths...),
			},
		}
	}

	return attributes
}

func (m *MonitorConfigScheduleIntervalResourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"year":   types.Int64Type,
		"month":  types.Int64Type,
		"week":   types.Int64Type,
		"day":    types.Int64Type,
		"hour":   types.Int64Type,
		"minute": types.Int64Type,
	}
}

func parseMonitorConfigScheduleInterval(m apiclient.MonitorConfigScheduleInterval) (MonitorConfigScheduleIntervalResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	rm := MonitorConfigScheduleIntervalResourceModel{}

	if len(m) != 2 {
		diags.AddError("Invalid schedule", "Invalid schedule")
		return rm, diags
	}

	var number int64
	number, ok := m[0].(int64)
	if !ok {
		diags.AddError("Invalid schedule", "Invalid schedule")
		return rm, diags
	}

	var unit string
	unit, ok = m[1].(string)
	if !ok {
		diags.AddError("Invalid schedule", "Invalid schedule")
		return rm, diags
	}

	switch unit {
	case "year":
		rm.Year = types.Int64Value(number)
	case "month":
		rm.Month = types.Int64Value(number)
	case "week":
		rm.Week = types.Int64Value(number)
	case "day":
		rm.Day = types.Int64Value(number)
	case "hour":
		rm.Hour = types.Int64Value(number)
	case "minute":
		rm.Minute = types.Int64Value(number)
	default:
		diags.AddError("Invalid schedule", "Invalid schedule")
	}

	return rm, diags
}
