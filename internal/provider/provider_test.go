package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/pkg/must"
	"github.com/jianyuan/terraform-provider-sentry/sentry"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	acctest.ProviderName: func() (tfprotov6.ProviderServer, error) {
		ctx := context.Background()

		upgradedSdkProvider := must.Get(tf5to6server.UpgradeServer(
			context.Background(),
			sentry.NewProvider(acctest.ProviderVersion)().GRPCProvider,
		))
		providers := []func() tfprotov6.ProviderServer{
			providerserver.NewProtocol6(New(acctest.ProviderVersion)()),
			func() tfprotov6.ProviderServer {
				return upgradedSdkProvider
			},
		}

		muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

		if err != nil {
			return nil, err
		}

		return muxServer.ProviderServer(), nil
	},
}
