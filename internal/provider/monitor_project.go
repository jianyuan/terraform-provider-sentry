package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

// fillMonitorProject sets a missing project (e.g. after import) to the slug of the project
// the monitor belongs to. The API only returns the numeric project ID, but configurations
// usually reference the slug, and a project change forces replacement.
func fillMonitorProject(
	ctx context.Context,
	apiClient *apiclient.ClientWithResponses,
	organization string,
	project *supertypes.StringValue,
	projectId string,
) (diags diag.Diagnostics) {
	if project.IsKnown() {
		return
	}

	if projectId == "" {
		diags.AddError("Client Error", "Unable to read project, monitor has no project ID")
		return
	}

	slug, err := sentryclient.GetProjectSlug(ctx, apiClient, organization, projectId)
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Unable to read project, got error: %s", err))
		return
	}

	project.Set(slug)
	return
}
