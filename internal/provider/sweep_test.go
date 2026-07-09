package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mzglinski/terraform-provider-sentry/internal/acctest"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	acctest.SetupShared(ctx)
	defer acctest.TeardownShared(ctx)

	resource.TestMain(m)
}
