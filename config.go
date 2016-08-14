package main

import "net/url"

type Config struct {
	Token   string
	BaseURL string
}

func (c *Config) Client() (interface{}, error) {
	if c.BaseURL != "" {
		_, err := url.Parse(c.BaseURL)
		if err != nil {
			return nil, err
		}
	}

	cl := NewClient(nil, c.Token, c.BaseURL)

	return cl, nil
}
