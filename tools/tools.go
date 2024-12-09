//go:build tools

package tools

import (
	// Documentation generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
	// API client generation
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)
