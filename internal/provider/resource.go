package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/providerdata"
)

type baseResource struct {
	client    *sentry.Client
	apiClient *apiclient.ClientWithResponses
}

func (r *baseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData := req.ProviderData.(*providerdata.ProviderData)

	r.client = providerData.Client
	r.apiClient = providerData.ApiClient
}
