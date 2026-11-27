package apiclient

import (
	"context"
	"net/http"
)

// WithQueryParam adds a query parameter that the generated client does not expose.
func WithQueryParam(key, value string) RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()
		return nil
	}
}
