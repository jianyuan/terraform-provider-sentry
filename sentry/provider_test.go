package sentry

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/jianyuan/go-utils/must"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	acctest.ProviderName: func() (tfprotov6.ProviderServer, error) {
		ctx := context.Background()

		upgradedSdkProvider := must.Get(tf5to6server.UpgradeServer(
			context.Background(),
			NewProvider(acctest.ProviderVersion)().GRPCProvider,
		))
		providers := []func() tfprotov6.ProviderServer{
			providerserver.NewProtocol6(provider.New(acctest.ProviderVersion)()),
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

func TestProvider(t *testing.T) {
	if err := NewProvider("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
