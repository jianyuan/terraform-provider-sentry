package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mzglinski/terraform-provider-sentry/internal/apiclient"
)

func (d *ProjectErrorMonitorDataSource) read(ctx context.Context, data *ProjectErrorMonitorDataSourceModel) (diags diag.Diagnostics) {
	listParams := &apiclient.ListOrganizationMonitorsParams{}

	// Resolve project ID
	projectHttpResp, err := d.apiClient.GetOrganizationProjectWithResponse(ctx, data.Organization.ValueString(), data.Project.ValueString())
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Unable to read project, got error: %s", err))
		return
	} else if projectHttpResp.StatusCode() != http.StatusOK {
		diags.AddError("Client Error", fmt.Sprintf("Unable to read project, got status code %d: %s", projectHttpResp.StatusCode(), string(projectHttpResp.Body)))
		return
	} else if projectHttpResp.JSON200 == nil {
		diags.AddError("Client Error", "Unable to read project, got empty response body")
		return
	}

	listParams.Project = new(projectHttpResp.JSON200.Id)

	// Build query
	var queryParts []string
	queryParts = append(queryParts, fmt.Sprintf("type:%s", "error"))
	listParams.Query = new(strings.Join(queryParts, " "))

	listHttpResp, err := d.apiClient.ListOrganizationMonitorsWithResponse(ctx, data.Organization.ValueString(), listParams)
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Unable to read monitors, got error: %s", err))
		return
	} else if listHttpResp.StatusCode() != http.StatusOK {
		diags.AddError("Client Error", fmt.Sprintf("Unable to read monitors, got status code %d: %s", listHttpResp.StatusCode(), string(listHttpResp.Body)))
		return
	} else if listHttpResp.JSON200 == nil {
		diags.AddError("Client Error", "Unable to read monitors, got empty response body")
		return
	}

	if len(*listHttpResp.JSON200) == 0 {
		diags.AddError("Client Error", "Unable to read monitors, no monitors found")
		return
	} else if len(*listHttpResp.JSON200) > 1 && !data.First.IsKnown() && !data.First.ValueBool() {
		diags.AddError("Client Error", "Multiple monitors found, please set the `first` attribute to `true`.")
		return
	}

	monitor := (*listHttpResp.JSON200)[0]
	diags.Append(data.Fill(ctx, monitor)...)

	return
}

func (m *ProjectErrorMonitorDataSourceModel) fill(ctx context.Context, data apiclient.ProjectMonitor) (diags diag.Diagnostics) {
	if data.Owner.IsSpecified() && !data.Owner.IsNull() {
		inOwner, err := data.Owner.Get()
		if err != nil {
			diags.AddError("Invalid owner", err.Error())
			return
		}

		inOwnerValue, err := inOwner.ValueByDiscriminator()
		if err != nil {
			diags.AddError("Invalid owner", err.Error())
			return
		}

		outOwner := &ProjectErrorMonitorDataSourceModelOwner{}

		switch inOwnerValue := inOwnerValue.(type) {
		case apiclient.ProjectMonitorOwnerUser:
			outOwner.UserId.Set(inOwnerValue.Id)
			diags.Append(m.Owner.Set(ctx, outOwner)...)
		case apiclient.ProjectMonitorOwnerTeam:
			outOwner.TeamId.Set(inOwnerValue.Id)
			diags.Append(m.Owner.Set(ctx, outOwner)...)
		default:
			m.Owner.SetNull(ctx)
		}
	} else {
		m.Owner.SetNull(ctx)
	}

	return
}
