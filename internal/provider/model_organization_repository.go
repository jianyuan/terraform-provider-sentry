package provider

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

type OrganizationRepositoryModel struct {
	Id              types.String `tfsdk:"id"`
	Organization    types.String `tfsdk:"organization"`
	IntegrationType types.String `tfsdk:"integration_type"`
	IntegrationId   types.String `tfsdk:"integration_id"`
	Identifier      types.String `tfsdk:"identifier"`
}

func (m *OrganizationRepositoryModel) Fill(organization string, repo sentry.OrganizationRepository) error {
	m.Id = types.StringValue(repo.ID)
	m.Organization = types.StringValue(organization)
	m.IntegrationType = types.StringValue(strings.TrimPrefix(repo.Provider.ID, "integrations:"))
	m.IntegrationId = types.StringValue(repo.IntegrationId)

	var identifierStr string
	var identifierNum json.Number
	if err := json.Unmarshal(repo.ExternalSlug, &identifierStr); err == nil {
		m.Identifier = types.StringValue(identifierStr)
	} else if err := json.Unmarshal(repo.ExternalSlug, &identifierNum); err == nil {
		m.Identifier = types.StringValue(identifierNum.String())
	} else {
		return fmt.Errorf("failed to unmarshal identifier: %s", repo.ExternalSlug)
	}

	return nil
}
