package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/diagutils"
	intresource "github.com/jianyuan/terraform-provider-sentry/internal/resource"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
)

type IntegrationPagerDutyModel struct {
	Id                      types.String `tfsdk:"id"`
	Organization            types.String `tfsdk:"organization"`
	IntegrationId           types.String `tfsdk:"integration_id"`
	Service                 types.String `tfsdk:"service"`
	IntegrationKey          types.String `tfsdk:"integration_key"`
	IntegrationKeyWo        types.String `tfsdk:"integration_key_wo"`
	IntegrationKeyWoVersion types.Int64  `tfsdk:"integration_key_wo_version"`
}

func (m *IntegrationPagerDutyModel) Fill(ctx context.Context, item apiclient.OrganizationIntegrationPagerDutyServiceTableItem) (diags diag.Diagnostics) {
	m.Id = types.StringValue(item.Id.String())
	m.Service = types.StringValue(item.Service)
	m.IntegrationKey = types.StringValue(item.IntegrationKey)
	return
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
var _ resource.ResourceWithConfigValidators = &IntegrationPagerDuty{}

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
	resp.Schema = superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "Manage a PagerDuty service integration.\n\n" + integrationKeyWriteOnlyNote,
		},
		Attributes: superschema.Attributes{
			"id": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of this resource.",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"organization": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The organization of this resource.",
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"integration_id": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the PagerDuty integration. Source from the URL `https://<organization>.sentry.io/settings/integrations/pagerduty/<integration-id>/` or use the `sentry_organization_integration` data source.",
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"service": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the PagerDuty service.",
					Required:            true,
				},
			},
			"integration_key": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The integration key of the PagerDuty service. Exactly one of `integration_key` or `integration_key_wo` must be set.",
					Optional:            true,
					Sensitive:           true,
				},
			},
			"integration_key_wo": superschema.StringAttribute{
				Resource: new(writeOnlyStringAttribute("integration_key", "The integration key of the PagerDuty service. Exactly one of `integration_key` or `integration_key_wo` must be set.")),
			},
			"integration_key_wo_version": superschema.Int64Attribute{
				Resource: new(writeOnlyVersionAttribute("integration_key")),
			},
		},
	}.GetResource(ctx)
}

func (r *IntegrationPagerDuty) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("integration_key"),
			path.MatchRoot("integration_key_wo"),
		),
	}
}

func (r *IntegrationPagerDuty) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IntegrationPagerDutyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrationKey, diags := resolveWriteOnlySecretFromConfig(ctx, req.Config, resp.Private, "integration_key", data.IntegrationKey, data.IntegrationKeyWoVersion, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if integrationKey.Value == nil {
		resp.Diagnostics.AddAttributeError(path.Root("integration_key"), "Missing integration key", "One of `integration_key` or `integration_key_wo` must be set.")
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

	specificIntegration, err := integration.AsOrganizationIntegrationPagerDuty()
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("unmarshal", err))
		return
	}

	idsSeen := map[json.Number]struct{}{}
	for _, item := range specificIntegration.ConfigData.ServiceTable {
		idsSeen[item.Id] = struct{}{}
	}

	specificIntegration.ConfigData.ServiceTable = append(specificIntegration.ConfigData.ServiceTable, apiclient.OrganizationIntegrationPagerDutyServiceTableItem{
		Service:        data.Service.ValueString(),
		IntegrationKey: *integrationKey.Value,
		Id:             json.Number("0"),
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

	specificIntegration, err = integration.AsOrganizationIntegrationPagerDuty()
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("unmarshal", err))
		return
	}

	var found *apiclient.OrganizationIntegrationPagerDutyServiceTableItem
	for _, item := range specificIntegration.ConfigData.ServiceTable {
		if item.Service == data.Service.ValueString() && item.IntegrationKey == *integrationKey.Value {
			if _, ok := idsSeen[item.Id]; !ok {
				found = &item
				break
			}
		}
	}

	if found == nil {
		resp.Diagnostics.Append(diagutils.NewClientError("create", fmt.Errorf("service table item not found: %s", data.IntegrationId.ValueString())))
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *found)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if integrationKey.InUse {
		data.IntegrationKey = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	resp.Diagnostics.Append(setWriteOnlyInUse(ctx, resp.Private, "integration_key", integrationKey.InUse)...)
}

func (r *IntegrationPagerDuty) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IntegrationPagerDutyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GetOrganizationIntegrationWithResponse(ctx, data.Organization.ValueString(), data.IntegrationId.ValueString())
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("read", err))
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	} else if httpResp.StatusCode() != http.StatusOK || httpResp.JSON200 == nil {
		resp.Diagnostics.Append(diagutils.NewClientStatusError("read", httpResp.StatusCode(), httpResp.Body))
		return
	}

	integration := *httpResp.JSON200

	specificIntegration, err := integration.AsOrganizationIntegrationPagerDuty()
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("unmarshal", err))
		return
	}

	var found *apiclient.OrganizationIntegrationPagerDutyServiceTableItem
	for _, i := range specificIntegration.ConfigData.ServiceTable {
		if i.Id.String() == data.Id.ValueString() {
			found = &i
			break
		}
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *found)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Sentry returns the integration key in full. Keep it out of the state when
	// it is managed through integration_key_wo.
	inUse, diags := writeOnlyInUse(ctx, req.Private, "integration_key")
	resp.Diagnostics.Append(diags...)
	if inUse {
		data.IntegrationKey = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IntegrationPagerDuty) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IntegrationPagerDutyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var priorVersion types.Int64
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("integration_key_wo_version"), &priorVersion)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrationKey, diags := resolveWriteOnlySecretFromConfig(ctx, req.Config, req.Private, "integration_key", data.IntegrationKey, data.IntegrationKeyWoVersion, &priorVersion)
	resp.Diagnostics.Append(diags...)
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

	specificIntegration, err := integration.AsOrganizationIntegrationPagerDuty()
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("unmarshal", err))
		return
	}

	var found *apiclient.OrganizationIntegrationPagerDutyServiceTableItem
	for i, item := range specificIntegration.ConfigData.ServiceTable {
		if item.Id.String() == data.Id.ValueString() {
			found = &specificIntegration.ConfigData.ServiceTable[i]
			break
		}
	}

	if found == nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", fmt.Errorf("service table item not found: %s", data.IntegrationId.ValueString())))
		resp.State.RemoveResource(ctx)
		return
	}

	found.Service = data.Service.ValueString()
	// A nil value means the write-only key has not changed, so keep the one Sentry has.
	if integrationKey.Value != nil {
		found.IntegrationKey = *integrationKey.Value
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

	specificIntegration, err = integration.AsOrganizationIntegrationPagerDuty()
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("unmarshal", err))
		return
	}

	found = nil
	for _, item := range specificIntegration.ConfigData.ServiceTable {
		if item.Id.String() == data.Id.ValueString() {
			found = &item
			break
		}
	}

	if found == nil {
		resp.Diagnostics.Append(diagutils.NewClientError("update", fmt.Errorf("service table item not found: %s", data.IntegrationId.ValueString())))
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, *found)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if integrationKey.InUse {
		data.IntegrationKey = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	resp.Diagnostics.Append(setWriteOnlyInUse(ctx, resp.Private, "integration_key", integrationKey.InUse)...)
}

func (r *IntegrationPagerDuty) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IntegrationPagerDutyModel

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

	specificIntegration, err := integration.AsOrganizationIntegrationPagerDuty()
	if err != nil {
		resp.Diagnostics.Append(diagutils.NewClientError("unmarshal", err))
		return
	}

	var found bool
	for i, item := range specificIntegration.ConfigData.ServiceTable {
		if item.Id.String() == data.Id.ValueString() {
			specificIntegration.ConfigData.ServiceTable = append(specificIntegration.ConfigData.ServiceTable[:i], specificIntegration.ConfigData.ServiceTable[i+1:]...)
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

func (r *IntegrationPagerDuty) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	intresource.ImportState3PartPath("organization", "integration_id", "id")(ctx, req, resp)
}
