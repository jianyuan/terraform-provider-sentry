package main

import (
	"flag"

	"github.com/deste-org/terraform-provider-sentry/sentry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug: debugMode,

		ProviderFunc: sentry.NewProvider(version),
	}

	plugin.Serve(opts)
}
