package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/resourceid"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

func (r *OrganizationUserMappingResource) getCreateJSONRequestBody(ctx context.Context, data OrganizationUserMappingResourceModel) (*apiclient.CreateOrganizationExternalUserJSONRequestBody, diag.Diagnostics) {
	var diags diag.Diagnostics

	externalActorID, err := strconv.Atoi(data.ExternalId.ValueString())
	if err != nil {
		diags.AddError("Invalid Attribute", fmt.Sprintf("Unable to convert external_id to integer: %s", err))
		return nil, diags
	}

	body := apiclient.CreateOrganizationExternalUserJSONRequestBody{
		UserId:        int(data.UserId.ValueInt64()),
		ExternalName:  data.ExternalName.ValueString(),
		Provider:      data.ExternalProvider.ValueString(),
		IntegrationId: int(data.IntegrationId.ValueInt64()),
		Id:            externalActorID,
	}

	return &body, diags
}

func (r *OrganizationUserMappingResource) getUpdateJSONRequestBody(ctx context.Context, data OrganizationUserMappingResourceModel) (*apiclient.UpdateOrganizationExternalUserJSONRequestBody, diag.Diagnostics) {
	createBody, diags := r.getCreateJSONRequestBody(ctx, data)
	if diags.HasError() || createBody == nil {
		return nil, diags
	}
	body := apiclient.UpdateOrganizationExternalUserJSONRequestBody(*createBody)
	return &body, diags
}

func (r *OrganizationUserMappingResource) read(ctx context.Context, data *OrganizationUserMappingResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// No GET API for external user mappings; keep prior state.
	if data.Id.IsNull() || data.Id.ValueString() == "" {
		if !data.Organization.IsNull() && !data.InternalId.IsNull() {
			id, err := resourceid.BuildPath2(data.Organization.ValueString(), data.InternalId.ValueString())
			if err != nil {
				diags.AddError("Invalid ID", err.Error())
				return diags
			}
			data.Id = supertypes.NewStringValue(id)
		}
	}
	return diags
}

func (m *OrganizationUserMappingResourceModel) Fill(ctx context.Context, data apiclient.ExternalUser) (diags diag.Diagnostics) {
	m.InternalId = supertypes.NewStringValue(data.Id)
	if !m.Organization.IsNull() && !m.Organization.IsUnknown() {
		id, err := resourceid.BuildPath2(m.Organization.ValueString(), data.Id)
		if err != nil {
			diags.AddError("Invalid ID", err.Error())
		} else {
			m.Id = supertypes.NewStringValue(id)
		}
	} else {
		m.Id = supertypes.NewStringValue(data.Id)
	}

	m.ExternalName = supertypes.NewStringValue(data.ExternalName)
	m.ExternalProvider = supertypes.NewStringValue(data.Provider)

	if userID, err := strconv.ParseInt(data.UserId, 10, 64); err == nil {
		m.UserId = supertypes.NewInt64Value(userID)
	} else {
		diags.AddError("Client Error", fmt.Sprintf("Unable to parse userId %q: %s", data.UserId, err))
	}
	if integrationID, err := strconv.ParseInt(data.IntegrationId, 10, 64); err == nil {
		m.IntegrationId = supertypes.NewInt64Value(integrationID)
	} else {
		diags.AddError("Client Error", fmt.Sprintf("Unable to parse integrationId %q: %s", data.IntegrationId, err))
	}

	return diags
}
