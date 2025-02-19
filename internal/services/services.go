package services

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/providerdata"
)

type BaseResource struct {
	Client    *sentry.Client
	ApiClient *apiclient.ClientWithResponses
}

func (r *BaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData := req.ProviderData.(*providerdata.ProviderData)

	r.Client = providerData.Client
	r.ApiClient = providerData.ApiClient
}

type BaseDataSource struct {
	Client    *sentry.Client
	ApiClient *apiclient.ClientWithResponses
}

func (d *BaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData := req.ProviderData.(*providerdata.ProviderData)

	d.Client = providerData.Client
	d.ApiClient = providerData.ApiClient
}
