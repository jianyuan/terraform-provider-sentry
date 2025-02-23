package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

type IntegrationOpsgenieModel struct {
	Id             types.String `tfsdk:"id"`
	Organization   types.String `tfsdk:"organization"`
	IntegrationId  types.String `tfsdk:"integration_id"`
	Team           types.String `tfsdk:"team"`
	IntegrationKey types.String `tfsdk:"integration_key"`
}

func (m *IntegrationOpsgenieModel) Fill(ctx context.Context, item apiclient.OrganizationIntegrationOpsgenieTeamTableItem) (diags diag.Diagnostics) {
	m.Id = types.StringValue(item.Id)
	m.Team = types.StringValue(item.Team)
	m.IntegrationKey = types.StringValue(item.IntegrationKey)
	return
}

var _ resource.Resource = &IntegrationOpsgenie{}
var _ resource.ResourceWithConfigure = &IntegrationOpsgenie{}
var _ resource.ResourceWithImportState = &IntegrationOpsgenie{}

func NewIntegrationOpsgenie() resource.Resource {
	return &IntegrationOpsgenie{}
}

type IntegrationOpsgenie struct {
	baseResource
}

func (r *IntegrationOpsgenie) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_opsgenie"
}

func (r *IntegrationOpsgenie) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage an Opsgenie team integration.",

		Attributes: map[string]schema.Attribute{
			"id":           ResourceIdAttribute(),
			"organization": ResourceOrganizationAttribute(),
			"integration_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Opsgenie integration. Source from the URL `https://<organization>.sentry.io/settings/integrations/opsgenie/<integration-id>/` or use the `sentry_organization_integration` data source.",
				Required:            true,
			},
			"team": schema.StringAttribute{
				MarkdownDescription: "The name of the Opsgenie team. In Sentry, this is called Label.",
				Required:            true,
			},
			"integration_key": schema.StringAttribute{
				MarkdownDescription: "The integration key of the Opsgenie service.",
				Required:            true,
			},
		},
	}
}

func (r *IntegrationOpsgenie) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IntegrationOpsgenieModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getHttpResp, err := r.apiClient.GetOrganizationIntegrationWithResponse(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if getHttpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("integration"))
		return
	} else if getHttpResp.StatusCode() != http.StatusOK || getHttpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", getHttpResp.StatusCode(), getHttpResp.Body))
		return
	}

	integration := *getHttpResp.JSON200

	specificIntegration, err := integration.AsOrganizationIntegrationOpsgenie()
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("unmarshal", err))
		return
	}

	idsSeen := map[string]struct{}{}
	for _, item := range specificIntegration.ConfigData.TeamTable {
		idsSeen[item.Id] = struct{}{}
	}

	specificIntegration.ConfigData.TeamTable = append(specificIntegration.ConfigData.TeamTable, apiclient.OrganizationIntegrationOpsgenieTeamTableItem{
		Team:           data.Team.ValueString(),
		IntegrationKey: data.IntegrationKey.ValueString(),
		Id:             "",
	})

	configDataJSON, err := json.Marshal(specificIntegration.ConfigData)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("marshal", err))
		return
	}

	updateHttpResp, err := r.apiClient.UpdateOrganizationIntegrationWithBodyWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.IntegrationId.ValueString(),
		"application/json",
		bytes.NewReader(configDataJSON),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	} else if updateHttpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("update", updateHttpResp.StatusCode(), updateHttpResp.Body))
		return
	}

	getHttpResp, err = r.apiClient.GetOrganizationIntegrationWithResponse(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if getHttpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("integration"))
		return
	} else if getHttpResp.StatusCode() != http.StatusOK || getHttpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", getHttpResp.StatusCode(), getHttpResp.Body))
		return
	}

	integration = *getHttpResp.JSON200

	specificIntegration, err = integration.AsOrganizationIntegrationOpsgenie()
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("unmarshal", err))
		return
	}

	var found *apiclient.OrganizationIntegrationOpsgenieTeamTableItem
	for _, item := range specificIntegration.ConfigData.TeamTable {
		if item.Team == data.Team.ValueString() && item.IntegrationKey == data.IntegrationKey.ValueString() {
			if _, ok := idsSeen[item.Id]; !ok {
				found = &item
				break
			}
		}
	}

	if found == nil {
		resp.Diagnostics.Append(diagutils.NewClientError("create", fmt.Errorf("team table item not found: %s", data.IntegrationId.ValueString())))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *found)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationOpsgenie) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IntegrationOpsgenieModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GetOrganizationIntegrationWithResponse(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("integration"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	integration := *httpResp.JSON200

	specificIntegration, err := integration.AsOrganizationIntegrationOpsgenie()
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("unmarshal", err))
		return
	}

	var found *apiclient.OrganizationIntegrationOpsgenieTeamTableItem
	for _, item := range specificIntegration.ConfigData.TeamTable {
		if item.Id == data.Id.ValueString() {
			found = &item
			break
		}
	}

	if found == nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", fmt.Errorf("team table item not found: %s", data.IntegrationId.ValueString())))
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *found)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationOpsgenie) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IntegrationOpsgenieModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GetOrganizationIntegrationWithResponse(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("integration"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	integration := *httpResp.JSON200

	specificIntegration, err := integration.AsOrganizationIntegrationOpsgenie()
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("unmarshal", err))
		return
	}

	var found *apiclient.OrganizationIntegrationOpsgenieTeamTableItem
	for i, item := range specificIntegration.ConfigData.TeamTable {
		if item.Id == data.Id.ValueString() {
			found = &specificIntegration.ConfigData.TeamTable[i]
			break
		}
	}

	if found == nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", fmt.Errorf("team table item not found: %s", data.IntegrationId.ValueString())))
		resp.State.RemoveResource(ctx)
		return
	}

	found.Team = data.Team.ValueString()
	found.IntegrationKey = data.IntegrationKey.ValueString()

	configDataJSON, err := json.Marshal(specificIntegration.ConfigData)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("marshal", err))
		return
	}

	updateHttpResp, err := r.apiClient.UpdateOrganizationIntegrationWithBodyWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.IntegrationId.ValueString(),
		"application/json",
		bytes.NewReader(configDataJSON),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("integration"))
		resp.State.RemoveResource(ctx)
		return
	} else if updateHttpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("update", updateHttpResp.StatusCode(), updateHttpResp.Body))
		return
	}

	httpResp, err = r.apiClient.GetOrganizationIntegrationWithResponse(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.Append(diagutils.NewNotFoundError("integration"))
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	integration = *httpResp.JSON200

	specificIntegration, err = integration.AsOrganizationIntegrationOpsgenie()
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("unmarshal", err))
		return
	}

	found = nil
	for _, item := range specificIntegration.ConfigData.TeamTable {
		if item.Id == data.Id.ValueString() {
			found = &item
			break
		}
	}

	if found == nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", fmt.Errorf("team table item not found: %s", data.IntegrationId.ValueString())))
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *found)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationOpsgenie) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IntegrationOpsgenieModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GetOrganizationIntegrationWithResponse(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	integration := *httpResp.JSON200

	specificIntegration, err := integration.AsOrganizationIntegrationOpsgenie()
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("unmarshal", err))
		return
	}

	var found bool
	for i, item := range specificIntegration.ConfigData.TeamTable {
		if item.Id == data.Id.ValueString() {
			specificIntegration.ConfigData.TeamTable = append(specificIntegration.ConfigData.TeamTable[:i], specificIntegration.ConfigData.TeamTable[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return
	}

	configDataJSON, err := json.Marshal(specificIntegration.ConfigData)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("marshal", err))
		return
	}

	updateHttpResp, err := r.apiClient.UpdateOrganizationIntegrationWithBodyWithResponse(
		ctx,
		data.Organization.ValueString(),
		data.IntegrationId.ValueString(),
		"application/json",
		bytes.NewReader(configDataJSON),
	)
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("delete", err))
		return
	} else if updateHttpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("delete", updateHttpResp.StatusCode(), updateHttpResp.Body))
		return
	}
}

func (r *IntegrationOpsgenie) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tfutils.ImportStateThreePartId(ctx, "organization", "integration_id", req, resp)
}
