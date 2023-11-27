package sentry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ProjectFilter represents inbounding filters applied to a project.
type ProjectFilter struct {
	ID     string          `json:"id"`
	Active json.RawMessage `json:"active"`
}

// ProjectOwnershipService provides methods for accessing Sentry project
// filters API endpoints.
type ProjectFilterService service

// Get the filters.
func (s *ProjectFilterService) Get(ctx context.Context, organizationSlug string, projectSlug string) ([]*ProjectFilter, *Response, error) {
	url := fmt.Sprintf("0/projects/%v/%v/filters/", organizationSlug, projectSlug)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	var filters []*ProjectFilter
	resp, err := s.client.Do(ctx, req, &filters)
	if err != nil {
		return nil, resp, err
	}

	return filters, resp, nil
}

// FilterConfig represents configuration for project filter
type FilterConfig struct {
	BrowserExtension bool
	LegacyBrowsers   []string
}

// GetFilterConfig retrieves filter configuration.
func (s *ProjectFilterService) GetFilterConfig(ctx context.Context, organizationSlug string, projectSlug string) (*FilterConfig, *Response, error) {
	filters, resp, err := s.Get(ctx, organizationSlug, projectSlug)
	if err != nil {
		return nil, resp, err
	}

	var filterConfig FilterConfig

	for _, filter := range filters {
		switch filter.ID {
		case "browser-extensions":
			if string(filter.Active) == "true" {
				filterConfig.BrowserExtension = true
			}

		case "legacy-browsers":
			if string(filter.Active) != "false" {
				err = json.Unmarshal(filter.Active, &filterConfig.LegacyBrowsers)
				if err != nil {
					return nil, resp, err
				}
			}
		}
	}

	return &filterConfig, resp, err
}

// BrowserExtensionParams defines parameters for browser extension request
type BrowserExtensionParams struct {
	Active bool `json:"active"`
}

// UpdateBrowserExtensions updates configuration for browser extension filter
func (s *ProjectFilterService) UpdateBrowserExtensions(ctx context.Context, organizationSlug string, projectSlug string, active bool) (*Response, error) {
	url := fmt.Sprintf("0/projects/%v/%v/filters/browser-extensions/", organizationSlug, projectSlug)
	params := BrowserExtensionParams{active}
	req, err := s.client.NewRequest(http.MethodPut, url, params)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}

// LegacyBrowserParams defines parameters for legacy browser request
type LegacyBrowserParams struct {
	Browsers []string `json:"subfilters"`
}

// UpdateLegacyBrowser updates configuration for legacy browser filters
func (s *ProjectFilterService) UpdateLegacyBrowser(ctx context.Context, organizationSlug string, projectSlug string, browsers []string) (*Response, error) {
	url := fmt.Sprintf("0/projects/%v/%v/filters/legacy-browsers/", organizationSlug, projectSlug)
	params := LegacyBrowserParams{browsers}

	req, err := s.client.NewRequest(http.MethodPut, url, params)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
