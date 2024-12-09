package sentryclient

import (
	"net/http"
	"sync"

	"github.com/jianyuan/go-sentry/v2/sentry"
	"golang.org/x/sync/semaphore"
)

func NewSemaphoreRoundTripper(delegate http.RoundTripper) http.RoundTripper {
	if delegate == nil {
		delegate = http.DefaultTransport
	}

	return &SemaphoreTransport{
		delegate: delegate,
	}
}

type SemaphoreTransport struct {
	delegate http.RoundTripper

	mu sync.RWMutex
	w  *semaphore.Weighted
}

func (t *SemaphoreTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mu.RLock()
	if t.w == nil {
		t.mu.RUnlock()
		t.mu.Lock()
		resp, err := t.delegate.RoundTrip(req)
		if resp != nil {
			rate := sentry.ParseRate(resp)
			if rate.ConcurrentLimit > 0 {
				t.w = semaphore.NewWeighted(int64(rate.ConcurrentLimit))
			}
		}
		t.mu.Unlock()
		return resp, err
	}
	t.mu.RUnlock()

	ctx := req.Context()
	if err := t.w.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer t.w.Release(1)

	return t.delegate.RoundTrip(req)
}
