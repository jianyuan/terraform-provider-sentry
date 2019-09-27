package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/jianyuan/terraform-provider-sentry/sentry"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sentry.Provider,
	})
}
