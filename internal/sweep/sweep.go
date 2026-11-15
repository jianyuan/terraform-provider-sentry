package sweep

import (
	"context"
	"errors"

	"github.com/jianyuan/terraform-provider-sentry/internal/providerdata"
	"golang.org/x/sync/errgroup"
)

// maxConcurrency keeps deletes well under Sentry's per-endpoint rate limit
// (e.g. 40 requests per second), so retries on 429 aren't exhausted.
const maxConcurrency = 10

type Sweepable interface {
	Delete(ctx context.Context) error
}

func Sweep(ctx context.Context, sweepables []Sweepable) error {
	errs := make([]error, len(sweepables))

	var g errgroup.Group
	g.SetLimit(maxConcurrency)
	for i, sweepable := range sweepables {
		g.Go(func() error {
			errs[i] = sweepable.Delete(ctx)
			return nil
		})
	}
	_ = g.Wait()

	return errors.Join(errs...)
}

type SweeperFn func(ctx context.Context, pd *providerdata.ProviderData) ([]Sweepable, error)
