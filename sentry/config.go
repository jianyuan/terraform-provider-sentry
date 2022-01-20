package sentry

import (
	"net/url"

	"github.com/jianyuan/go-sentry/sentry"
	"github.com/jianyuan/terraform-provider-sentry/logging"
)

// Config is the configuration structure used to instantiate the Sentry
// provider.
type Config struct {
	Token   string
	BaseURL string
}

// Client to connect to Sentry.
func (c *Config) Client() (interface{}, error) {
	var baseURL *url.URL
	var err error

	if c.BaseURL != "" {
		baseURL, err = url.Parse(c.BaseURL)
		logging.Errorf("Parsing base url %s", c.BaseURL)
		if err != nil {
			return nil, err
		}
	} else {
		logging.Warning("No base URL was set for the Sentry client")
	}

	logging.Info("Instantiating Sentry client...")
	cl := sentry.NewClient(nil, baseURL, c.Token)

	return cl, nil
}
