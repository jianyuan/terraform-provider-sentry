package sweep

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jianyuan/terraform-provider-sentry/internal/acctest"
)

func Register(name string, fn SweeperFn, dependencies ...string) {
	resource.AddTestSweepers(name, &resource.Sweeper{
		Name: name,
		F: func(region string) error {
			ctx := context.Background()

			sweepables, err := fn(ctx, acctest.SharedProviderData)
			if err != nil {
				return fmt.Errorf("failed to collect %q: %w", name, err)
			}

			err = Sweep(ctx, sweepables)
			if err != nil {
				return fmt.Errorf("failed to sweep %q: %w", name, err)
			}

			return nil
		},
		Dependencies: dependencies,
	})
}
