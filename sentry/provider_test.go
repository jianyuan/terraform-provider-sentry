package sentry

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
	"github.com/jianyuan/terraform-provider-sentry/internal/provider"
)

var testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
	acctest.ProviderName: func() (tfprotov5.ProviderServer, error) {
		ctx := context.Background()
		providers := []func() tfprotov5.ProviderServer{
			providerserver.NewProtocol5(provider.New(acctest.ProviderVersion)()),
			NewProvider(acctest.ProviderVersion)().GRPCProvider,
		}

		muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

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
