package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		BaseURL:   baseUrl,
	}
	client, err := config.Client(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to create Sentry client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SentryProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectInboundDataFilterResource,
		NewProjectSpikeProtectionResource,
		NewProjectSymbolSourcesResource,
		NewTeamMemberResource,
	}
}

func (p *SentryProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SentryProvider{
			version: version,
		}
	}
}
