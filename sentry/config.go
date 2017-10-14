package sentry

import (
	"log"
	"net/url"

	"github.com/jianyuan/go-sentry/sentry"
)

// Config is the configuration structure used to instantiate the Sentry
// provider.
type Config struct {
	Token   string
	BaseURL string
}

func (c *Config) Client() (interface{}, error) {
	var baseURL *url.URL
	var err error

	if c.BaseURL != "" {
		baseURL, err = url.Parse(c.BaseURL)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("[INFO] Instantiating Sentry client...")
	cl := sentry.NewClient(nil, baseURL, c.Token)

	return cl, nil
}
