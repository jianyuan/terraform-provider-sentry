package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &IntegrationOpsgenie{}
var _ resource.ResourceWithConfigure = &IntegrationOpsgenie{}
var _ resource.ResourceWithImportState = &IntegrationOpsgenie{}

func NewIntegrationOpsgenie() resource.Resource {
	return &IntegrationOpsgenie{}
}

type IntegrationOpsgenie struct {
	baseResource
}

type IntegrationOpsgenieModel struct {
	Id             types.String `tfsdk:"id"`
	Organization   types.String `tfsdk:"organization"`
	IntegrationId  types.String `tfsdk:"integration_id"`
	Team           types.String `tfsdk:"team"`
	IntegrationKey types.String `tfsdk:"integration_key"`
}

func (m *IntegrationOpsgenieModel) Fill(organization string, integrationId string, item IntegrationOpsgenieConfigDataTeamTableItem) error {
	m.Id = types.StringValue(item.Id)
	m.Organization = types.StringValue(organization)
	m.IntegrationId = types.StringValue(integrationId)
	m.Team = types.StringValue(item.Team)
	m.IntegrationKey = types.StringValue(item.IntegrationKey)

	return nil
}

type IntegrationOpsgenieConfigDataTeamTableItem struct {
	Team           string `json:"team"`
	IntegrationKey string `json:"integration_key"`
	Id             string `json:"id"`
}

type IntegrationOpsgenieConfigData struct {
	TeamTable []IntegrationOpsgenieConfigDataTeamTableItem `json:"team_table"`
}

func (r *IntegrationOpsgenie) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_opsgenie"
}

func (r *IntegrationOpsgenie) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage an Opsgenie team integration.",

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

	integration, apiResp, err := r.client.OrganizationIntegrations.Get(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Not found: %s", err.Error()))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", err.Error()))
		return
	}

	var configData IntegrationOpsgenieConfigData
	if err := json.Unmarshal(integration.ConfigData, &configData); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unmarshal error: %s", err.Error()))
		return
	}

	idsSeen := map[string]struct{}{}
	for _, item := range configData.TeamTable {
		idsSeen[item.Id] = struct{}{}
	}

	configData.TeamTable = append(configData.TeamTable, IntegrationOpsgenieConfigDataTeamTableItem{
		Team:           data.Team.ValueString(),
		IntegrationKey: data.IntegrationKey.ValueString(),
		Id:             "",
	})

	configDataJSON, err := json.Marshal(configData)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Marshal error: %s", err.Error()))
		return
	}

	params := json.RawMessage(configDataJSON)
	_, err = r.client.OrganizationIntegrations.UpdateConfig(
		ctx,
		data.Organization.ValueString(),
		data.IntegrationId.ValueString(),
		&params,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Create error: %s", err.Error()))
		return
	}

	integration, apiResp, err = r.client.OrganizationIntegrations.Get(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Not found: %s", err.Error()))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", err.Error()))
		return
	}

	if err := json.Unmarshal(integration.ConfigData, &configData); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unmarshal error: %s", err.Error()))
		return
	}

	var found *IntegrationOpsgenieConfigDataTeamTableItem
	for _, item := range configData.TeamTable {
		if item.Team == data.Team.ValueString() && item.IntegrationKey == data.IntegrationKey.ValueString() {
			if _, ok := idsSeen[item.Id]; !ok {
				found = &item
				break
			}
		}
	}

	if found == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Team table item not found: %s", data.Team.ValueString()))
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.IntegrationId.ValueString(), configData.TeamTable[len(configData.TeamTable)-1]); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Fill error: %s", err.Error()))
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

	integration, apiResp, err := r.client.OrganizationIntegrations.Get(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Not found: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", err.Error()))
		return
	}

	var configData IntegrationOpsgenieConfigData
	if err := json.Unmarshal(integration.ConfigData, &configData); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unmarshal error: %s", err.Error()))
		return
	}

	var found *IntegrationOpsgenieConfigDataTeamTableItem
	for _, i := range configData.TeamTable {
		if i.Id == data.Id.ValueString() {
			found = &i
			break
		}
	}
	if found == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Team table item not found: %s", data.IntegrationId.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.IntegrationId.ValueString(), *found); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Fill error: %s", err.Error()))
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

	integration, apiResp, err := r.client.OrganizationIntegrations.Get(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", err.Error()))
		return
	}

	var configData IntegrationOpsgenieConfigData
	if err := json.Unmarshal(integration.ConfigData, &configData); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unmarshal error: %s", err.Error()))
		return
	}

	var found *IntegrationOpsgenieConfigDataTeamTableItem
	for i, item := range configData.TeamTable {
		if item.Id == data.Id.ValueString() {
			found = &configData.TeamTable[i]
			break
		}
	}

	if found == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Team table item not found: %s", data.IntegrationId.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	found.Team = data.Team.ValueString()
	found.IntegrationKey = data.IntegrationKey.ValueString()

	configDataJSON, err := json.Marshal(configData)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Marshal error: %s", err.Error()))
		return
	}

	params := json.RawMessage(configDataJSON)
	_, err = r.client.OrganizationIntegrations.UpdateConfig(
		ctx,
		data.Organization.ValueString(),
		data.IntegrationId.ValueString(),
		&params,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Update error: %s", err.Error()))
		return
	}

	integration, apiResp, err = r.client.OrganizationIntegrations.Get(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if apiResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Not found: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", err.Error()))
		return
	}

	if err := json.Unmarshal(integration.ConfigData, &configData); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unmarshal error: %s", err.Error()))
		return
	}

	found = nil
	for _, item := range configData.TeamTable {
		if item.Id == data.Id.ValueString() {
			found = &item
			break
		}
	}
	if found == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Team table item not found: %s", data.IntegrationId.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.IntegrationId.ValueString(), *found); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Fill error: %s", err.Error()))
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

	integration, apiResp, err := r.client.OrganizationIntegrations.Get(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Read error: %s", err.Error()))
		return
	}

	var configData IntegrationOpsgenieConfigData
	if err := json.Unmarshal(integration.ConfigData, &configData); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unmarshal error: %s", err.Error()))
		return
	}

	var found bool
	for i, item := range configData.TeamTable {
		if item.Id == data.Id.ValueString() {
			configData.TeamTable = append(configData.TeamTable[:i], configData.TeamTable[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return
	}

	configDataJSON, err := json.Marshal(configData)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Marshal error: %s", err.Error()))
		return
	}

	params := json.RawMessage(configDataJSON)
	_, err = r.client.OrganizationIntegrations.UpdateConfig(
		ctx,
		data.Organization.ValueString(),
		data.IntegrationId.ValueString(),
		&params,
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Delete error: %s", err.Error()))
		return
	}
}

func (r *IntegrationOpsgenie) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, integrationId, id, err := splitThreePartID(req.ID, "organization", "integration-id", "id")
	if err != nil {
		resp.Diagnostics.AddError("Import Error", fmt.Sprintf("Unable to import integration, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("organization"), organization,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("integration_id"), integrationId,
	)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("id"), id,
	)...)
}
