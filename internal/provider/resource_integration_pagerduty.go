package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
)

type IntegrationPagerDutyModel struct {
	Id             types.String `tfsdk:"id"`
	Organization   types.String `tfsdk:"organization"`
	IntegrationId  types.String `tfsdk:"integration_id"`
	Service        types.String `tfsdk:"service"`
	IntegrationKey types.String `tfsdk:"integration_key"`
}

func (m *IntegrationPagerDutyModel) Fill(organization string, integrationId string, item IntegrationPagerDutyConfigDataServiceTableItem) error {
	m.Id = types.StringValue(item.Id.String())
	m.Organization = types.StringValue(organization)
	m.IntegrationId = types.StringValue(integrationId)
	m.Service = types.StringValue(item.Service)
	m.IntegrationKey = types.StringValue(item.IntegrationKey)

	return nil
}

type IntegrationPagerDutyConfigDataServiceTableItem struct {
	Service        string      `json:"service"`
	IntegrationKey string      `json:"integration_key"`
	Id             json.Number `json:"id"`
}

type IntegrationPagerDutyConfigData struct {
	ServiceTable []IntegrationPagerDutyConfigDataServiceTableItem `json:"service_table"`
}

var _ resource.Resource = &IntegrationPagerDuty{}
var _ resource.ResourceWithConfigure = &IntegrationPagerDuty{}
var _ resource.ResourceWithImportState = &IntegrationPagerDuty{}

func NewIntegrationPagerDuty() resource.Resource {
	return &IntegrationPagerDuty{}
}

type IntegrationPagerDuty struct {
	baseResource
}

func (r *IntegrationPagerDuty) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_pagerduty"
}

func (r *IntegrationPagerDuty) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage a PagerDuty service integration.",

		Attributes: map[string]schema.Attribute{
			"id":           ResourceIdAttribute(),
			"organization": ResourceOrganizationAttribute(),
			"integration_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the PagerDuty integration. Source from the URL `https://<organization>.sentry.io/settings/integrations/pagerduty/<integration-id>/` or use the `sentry_organization_integration` data source.",
				Required:            true,
			},
			"service": schema.StringAttribute{
				MarkdownDescription: "The name of the PagerDuty service.",
				Required:            true,
			},
			"integration_key": schema.StringAttribute{
				MarkdownDescription: "The integration key of the PagerDuty service.",
				Required:            true,
			},
		},
	}
}

func (r *IntegrationPagerDuty) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IntegrationPagerDutyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, apiResp, err := r.client.OrganizationIntegrations.Get(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if apiResp.StatusCode == http.StatusNotFound {
		diagutils.AddNotFoundError(resp.Diagnostics, "integration")
		return
	}
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "read", err)
		return
	}

	var configData IntegrationPagerDutyConfigData
	if err := json.Unmarshal(integration.ConfigData, &configData); err != nil {
		diagutils.AddClientError(resp.Diagnostics, "unmarshal", err)
		return
	}

	idsSeen := map[json.Number]struct{}{}
	for _, item := range configData.ServiceTable {
		idsSeen[item.Id] = struct{}{}
	}

	configData.ServiceTable = append(configData.ServiceTable, IntegrationPagerDutyConfigDataServiceTableItem{
		Service:        data.Service.ValueString(),
		IntegrationKey: data.IntegrationKey.ValueString(),
		Id:             json.Number("0"),
	})

	configDataJSON, err := json.Marshal(configData)
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "marshal", err)
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
		diagutils.AddClientError(resp.Diagnostics, "create", err)
		return
	}

	integration, apiResp, err = r.client.OrganizationIntegrations.Get(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if apiResp.StatusCode == http.StatusNotFound {
		diagutils.AddNotFoundError(resp.Diagnostics, "integration")
		return
	}
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "read", err)
		return
	}

	if err := json.Unmarshal(integration.ConfigData, &configData); err != nil {
		diagutils.AddClientError(resp.Diagnostics, "unmarshal", err)
		return
	}

	var found *IntegrationPagerDutyConfigDataServiceTableItem
	for _, item := range configData.ServiceTable {
		if item.Service == data.Service.ValueString() && item.IntegrationKey == data.IntegrationKey.ValueString() {
			if _, ok := idsSeen[item.Id]; !ok {
				found = &item
				break
			}
		}
	}

	if found == nil {
		diagutils.AddClientError(resp.Diagnostics, "create", fmt.Errorf("service table item not found: %s", data.IntegrationId.ValueString()))
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.IntegrationId.ValueString(), configData.ServiceTable[len(configData.ServiceTable)-1]); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationPagerDuty) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IntegrationPagerDutyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, apiResp, err := r.client.OrganizationIntegrations.Get(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if apiResp.StatusCode == http.StatusNotFound {
		diagutils.AddNotFoundError(resp.Diagnostics, "integration")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "read", err)
		return
	}

	var configData IntegrationPagerDutyConfigData
	if err := json.Unmarshal(integration.ConfigData, &configData); err != nil {
		diagutils.AddClientError(resp.Diagnostics, "unmarshal", err)
		return
	}

	var found *IntegrationPagerDutyConfigDataServiceTableItem
	for _, i := range configData.ServiceTable {
		if i.Id.String() == data.Id.ValueString() {
			found = &i
			break
		}
	}
	if found == nil {
		diagutils.AddClientError(resp.Diagnostics, "read", fmt.Errorf("service table item not found: %s", data.IntegrationId.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.IntegrationId.ValueString(), *found); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationPagerDuty) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IntegrationPagerDutyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, apiResp, err := r.client.OrganizationIntegrations.Get(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "update", err)
		return
	}

	var configData IntegrationPagerDutyConfigData
	if err := json.Unmarshal(integration.ConfigData, &configData); err != nil {
		diagutils.AddClientError(resp.Diagnostics, "unmarshal", err)
		return
	}

	var found *IntegrationPagerDutyConfigDataServiceTableItem
	for i, item := range configData.ServiceTable {
		if item.Id.String() == data.Id.ValueString() {
			found = &configData.ServiceTable[i]
			break
		}
	}

	if found == nil {
		diagutils.AddClientError(resp.Diagnostics, "update", fmt.Errorf("service table item not found: %s", data.IntegrationId.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	found.Service = data.Service.ValueString()
	found.IntegrationKey = data.IntegrationKey.ValueString()

	configDataJSON, err := json.Marshal(configData)
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "marshal", err)
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
		diagutils.AddClientError(resp.Diagnostics, "update", err)
		return
	}

	integration, apiResp, err = r.client.OrganizationIntegrations.Get(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if apiResp.StatusCode == http.StatusNotFound {
		diagutils.AddNotFoundError(resp.Diagnostics, "integration")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "update", err)
		return
	}

	if err := json.Unmarshal(integration.ConfigData, &configData); err != nil {
		diagutils.AddClientError(resp.Diagnostics, "unmarshal", err)
		return
	}

	found = nil
	for _, item := range configData.ServiceTable {
		if item.Id.String() == data.Id.ValueString() {
			found = &item
			break
		}
	}
	if found == nil {
		diagutils.AddClientError(resp.Diagnostics, "update", fmt.Errorf("service table item not found: %s", data.IntegrationId.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	if err := data.Fill(data.Organization.ValueString(), data.IntegrationId.ValueString(), *found); err != nil {
		diagutils.AddFillError(resp.Diagnostics, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationPagerDuty) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IntegrationPagerDutyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, apiResp, err := r.client.OrganizationIntegrations.Get(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if apiResp.StatusCode == http.StatusNotFound {
		return
	}
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "delete", err)
		return
	}

	var configData IntegrationPagerDutyConfigData
	if err := json.Unmarshal(integration.ConfigData, &configData); err != nil {
		diagutils.AddClientError(resp.Diagnostics, "unmarshal", err)
		return
	}

	var found bool
	for i, item := range configData.ServiceTable {
		if item.Id.String() == data.Id.ValueString() {
			configData.ServiceTable = append(configData.ServiceTable[:i], configData.ServiceTable[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return
	}

	configDataJSON, err := json.Marshal(configData)
	if err != nil {
		diagutils.AddClientError(resp.Diagnostics, "marshal", err)
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
		diagutils.AddClientError(resp.Diagnostics, "delete", err)
		return
	}
}

func (r *IntegrationPagerDuty) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	organization, integrationId, id, err := splitThreePartID(req.ID, "organization", "integration-id", "id")
	if err != nil {
		diagutils.AddImportError(resp.Diagnostics, err)
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
