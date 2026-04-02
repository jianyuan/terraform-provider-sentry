package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type customFiller[T any] interface {
	fill(ctx context.Context, data T) (diags diag.Diagnostics)
}
