package sentryclient

import "net/http"

func NewBearerTokenRoundTripper(delegate http.RoundTripper, token string) http.RoundTripper {
	if delegate == nil {
		delegate = http.DefaultTransport
	}

	return &BearerTokenRoundTripper{
		delegate: delegate,
		token:    token,
	}
}

type BearerTokenRoundTripper struct {
	delegate http.RoundTripper
	token    string
}

func (t *BearerTokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.delegate.RoundTrip(req)
}
