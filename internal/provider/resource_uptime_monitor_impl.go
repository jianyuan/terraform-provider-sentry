package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentrydata"
	"github.com/oapi-codegen/nullable"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (r *UptimeMonitorResource) getCreateJSONRequestBody(ctx context.Context, data UptimeMonitorResourceModel) (*apiclient.CreateProjectMonitorJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	outDs := apiclient.ProjectMonitorDataSourceUptimeDomainFailure{
		Url:             data.Url.Get(),
		Method:          data.Method.Get(),
		Headers:         [][]string{},
		IntervalSeconds: data.IntervalSeconds.Get(),
		TimeoutMs:       data.TimeoutMs.Get(),
	}
	if data.Body.IsKnown() {
		outDs.Body = nullable.NewNullableWithValue(data.Body.Get())
	} else {
		outDs.Body = nullable.NewNullNullable[string]()
	}
	if data.Headers.IsKnown() {
		inHeaders := data.Headers.DiagsGet(ctx, diags)
		if diags.HasError() {
			return nil, diags
		}
		for _, inHeader := range inHeaders {
			outDs.Headers = append(outDs.Headers, []string{inHeader.Key.Get(), inHeader.Value.Get()})
		}
	}

	var outConfig apiclient.ProjectMonitorConfig
	if err := outConfig.FromProjectMonitorConfigUptimeDomainFailure(apiclient.ProjectMonitorConfigUptimeDomainFailure{
		Mode:              apiclient.ProjectMonitorConfigUptimeDomainFailureMode(sentrydata.UptimeMonitorModeNameToId["MANUAL"]),
		Environment:       data.Environment.Get(),
		RecoveryThreshold: data.RecoveryThreshold.Get(),
		DowntimeThreshold: data.DowntimeThreshold.Get(),
	}); err != nil {
		diags.AddError("Error marshalling JSON", err.Error())
		return nil, diags
	}

	out := apiclient.ProjectMonitorRequestUptimeDomainFailure{
		Name:      data.Name.Get(),
		ProjectId: data.Project.Get(),
		DataSources: []apiclient.ProjectMonitorDataSourceUptimeDomainFailure{
			outDs,
		},
		Config: &outConfig,
	}

	if data.Enabled.IsKnown() {
		out.Enabled = nullable.NewNullableWithValue(data.Enabled.Get())
	} else {
		out.Enabled = nullable.NewNullNullable[bool]()
	}

	if data.Description.IsKnown() {
		out.Description = nullable.NewNullableWithValue(data.Description.Get())
	} else {
		out.Description = nullable.NewNullNullable[string]()
	}

	if data.DefaultAssignee.IsKnown() {
		defaultAssignee := data.DefaultAssignee.MustGet(ctx)
		switch {
		case defaultAssignee.TeamId.IsKnown():
			out.Owner = nullable.NewNullableWithValue(fmt.Sprintf("team:%s", defaultAssignee.TeamId.Get()))
		case defaultAssignee.UserId.IsKnown():
			out.Owner = nullable.NewNullableWithValue(fmt.Sprintf("user:%s", defaultAssignee.UserId.Get()))
		default:
			out.Owner = nullable.NewNullNullable[string]()
		}
	} else {
		out.Owner = nullable.NewNullNullable[string]()
	}

	var req apiclient.CreateProjectMonitorJSONRequestBody
	if err := req.FromProjectMonitorRequestUptimeDomainFailure(out); err != nil {
		diags.AddError("Error marshalling JSON", err.Error())
		return nil, diags
	}
	return &req, nil
}

func (m *UptimeMonitorResourceModel) Fill(ctx context.Context, data apiclient.ProjectMonitor) (diags diag.Diagnostics) {
	m.Id = supertypes.NewStringValue(data.Id)
	m.Name = supertypes.NewStringValue(data.Name)
	if v, err := data.Description.Get(); err == nil {
		m.Description = supertypes.NewStringValueOrNull(v)
	} else {
		m.Description = supertypes.NewStringNull()
	}
	m.Enabled = supertypes.NewBoolValue(data.Enabled)

	if data.Owner.IsSpecified() && !data.Owner.IsNull() {
		ownerValue, err := data.Owner.MustGet().ValueByDiscriminator()
		if err != nil {
			diags.AddError("Invalid owner", err.Error())
			return
		}

		defaultAssignee := &UptimeMonitorResourceModelDefaultAssignee{}

		switch ownerValue := ownerValue.(type) {
		case apiclient.ProjectMonitorOwnerUser:
			defaultAssignee.UserId = supertypes.NewStringValue(ownerValue.Id)
			m.DefaultAssignee = supertypes.NewSingleNestedObjectValueOf(ctx, defaultAssignee)
		case apiclient.ProjectMonitorOwnerTeam:
			defaultAssignee.TeamId = supertypes.NewStringValue(ownerValue.Id)
			m.DefaultAssignee = supertypes.NewSingleNestedObjectValueOf(ctx, defaultAssignee)
		default:
			m.DefaultAssignee = supertypes.NewSingleNestedObjectValueOfNull[UptimeMonitorResourceModelDefaultAssignee](ctx)
		}
	} else {
		m.DefaultAssignee = supertypes.NewSingleNestedObjectValueOfNull[UptimeMonitorResourceModelDefaultAssignee](ctx)
	}

	if len(data.DataSources) != 1 {
		diags.AddError("Invalid data source", fmt.Sprintf("Expected 1 data source, got %d", len(data.DataSources)))
		return
	}

	dataSource, err := data.DataSources[0].AsProjectMonitorDataSourceWrapperUptimeSubscription()
	if err != nil {
		diags.AddError("Invalid data source", err.Error())
		return
	}

	m.Url = supertypes.NewStringValue(dataSource.QueryObj.Url)
	m.Method = supertypes.NewStringValue(dataSource.QueryObj.Method)
	if v, err := dataSource.QueryObj.Body.Get(); err == nil {
		m.Body = supertypes.NewStringValueOrNull(v)
	} else {
		m.Body = supertypes.NewStringNull()
	}

	var headers []UptimeMonitorResourceModelHeadersItem
	for _, header := range dataSource.QueryObj.Headers {
		if len(header) != 2 {
			diags.AddError("Invalid header", fmt.Sprintf("Expected 2 elements in header, got %d", len(header)))
			return
		}
		headers = append(headers, UptimeMonitorResourceModelHeadersItem{
			Key:   supertypes.NewStringValue(header[0]),
			Value: supertypes.NewStringValue(header[1]),
		})
	}
	m.Headers = supertypes.NewListNestedObjectValueOfValueSlice(ctx, headers)

	m.IntervalSeconds = supertypes.NewInt64Value(dataSource.QueryObj.IntervalSeconds)
	m.TimeoutMs = supertypes.NewInt64Value(dataSource.QueryObj.TimeoutMs)

	if config, err := data.Config.AsProjectMonitorConfigUptimeDomainFailure(); err == nil {
		m.Environment = supertypes.NewStringValue(config.Environment)
		m.RecoveryThreshold = supertypes.NewInt64Value(config.RecoveryThreshold)
		m.DowntimeThreshold = supertypes.NewInt64Value(config.DowntimeThreshold)
	} else {
		diags.AddError("Invalid config", err.Error())
		return
	}

	return
}
