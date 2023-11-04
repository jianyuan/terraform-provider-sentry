package sentry

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/jianyuan/terraform-provider-sentry/internal/provider"
)

var testOrganization = os.Getenv("SENTRY_TEST_ORGANIZATION")

var testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
	"sentry": func() (tfprotov5.ProviderServer, error) {
		ctx := context.Background()
		providers := []func() tfprotov5.ProviderServer{
			providerserver.NewProtocol5(provider.New("test")()),
			NewProvider("test")().GRPCProvider,
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

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SENTRY_AUTH_TOKEN"); v == "" {
		t.Fatal("SENTRY_AUTH_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("SENTRY_TEST_ORGANIZATION"); v == "" {
		t.Fatal("SENTRY_TEST_ORGANIZATION must be set for acceptance tests")
	}
}
