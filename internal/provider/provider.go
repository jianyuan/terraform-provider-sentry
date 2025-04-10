package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jianyuan/go-sentry/v2/sentry"
	"github.com/jianyuan/terraform-provider-sentry/internal/apiclient"
	"github.com/jianyuan/terraform-provider-sentry/internal/providerdata"
	"github.com/jianyuan/terraform-provider-sentry/internal/sentryclient"
)

var _ provider.Provider = &SentryProvider{}

// SentryProvider defines the provider implementation.
type SentryProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SentryProviderModel describes the provider data model.
type SentryProviderModel struct {
	Token   types.String `tfsdk:"token"`
	BaseUrl types.String `tfsdk:"base_url"`
}

func (p *SentryProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sentry"
	resp.Version = p.version
}

func (p *SentryProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				MarkdownDescription: "The authentication token used to connect to Sentry. The value can be sourced from the `SENTRY_AUTH_TOKEN` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "The target Sentry Base API URL in the format `https://[hostname]/api/`. The default value is `https://sentry.io/api/`. The value must be provided when working with Sentry On-Premise. The value can be sourced from the `SENTRY_BASE_URL` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *SentryProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SentryProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var token string
	if !data.Token.IsNull() {
		token = data.Token.ValueString()
	} else if v := os.Getenv("SENTRY_AUTH_TOKEN"); v != "" {
		token = v
	} else if v := os.Getenv("SENTRY_TOKEN"); v != "" {
		token = v
	}

	var baseUrl string
	if !data.BaseUrl.IsNull() {
		baseUrl = data.BaseUrl.ValueString()
	} else if v := os.Getenv("SENTRY_BASE_URL"); v != "" {
		baseUrl = v
	} else {
		baseUrl = "https://sentry.io/api/"
	}

	config := sentryclient.Config{
		UserAgent: fmt.Sprintf("Terraform/%s (+https://www.terraform.io) terraform-provider-sentry/%s", req.TerraformVersion, p.version),
		Token:     token,
	}

	httpClient := config.HttpClient(ctx)

	// Old Sentry client
	var client *sentry.Client
	var err error
	if baseUrl == "" {
		client = sentry.NewClient(httpClient)
	} else {
		client, err = sentry.NewOnPremiseClient(baseUrl, httpClient)
		if err != nil {
			resp.Diagnostics.AddError("failed to create Sentry client", err.Error())
			return
		}
	}

	// New Sentry client
	apiClient, err := apiclient.NewClientWithResponses(
		baseUrl,
		apiclient.WithHTTPClient(httpClient),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to create Sentry API client", err.Error())
		return
	}

	httpResp, err := apiClient.HealthCheckWithResponse(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to perform health check", err.Error())
		return
	} else if httpResp.StatusCode() == http.StatusNotFound {
		resp.Diagnostics.AddError("failed to perform health check", "Sentry API is not available, please check the base URL")
		return
	} else if httpResp.StatusCode() == http.StatusUnauthorized {
		resp.Diagnostics.AddError("failed to perform health check", "Sentry API is not available, Please check the authentication token")
		return
	} else if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("failed to perform health check", fmt.Errorf("unexpected status code: %d", httpResp.StatusCode()).Error())
		return
	}

	providerData := &providerdata.ProviderData{
		Client:    client,
		ApiClient: apiClient,
	}

	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

func (p *SentryProvider) Resources(ctx context.Context) []func() resource.Resource {
	// Please keep the resources sorted by name.
	return []func() resource.Resource{
		NewAllProjectsSpikeProtectionResource,
		NewClientKeyResource,
		NewIntegrationOpsgenie,
		NewIntegrationPagerDuty,
		NewIssueAlertResource,
		NewMonitorResource,
		NewNotificationActionResource,
		NewOrganizationRepositoryResource,
		NewProjectInboundDataFilterResource,
		NewProjectResource,
		NewProjectSpikeProtectionResource,
		NewProjectSymbolSourcesResource,
		NewTeamMemberResource,
	}
}

func (p *SentryProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	// Please keep the data sources sorted by name.
	return []func() datasource.DataSource{
		NewAllClientKeysDataSource,
		NewAllOrganizationMembersDataSource,
		NewAllProjectsDataSource,
		NewClientKeyDataSource,
		NewIssueAlertDataSource,
		NewOrganizationDataSource,
		NewOrganizationIntegrationDataSource,
		NewOrganizationMemberDataSource,
		NewProjectDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SentryProvider{
			version: version,
		}
	}
}
