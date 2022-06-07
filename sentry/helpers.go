package sentry

import (
	"github.com/jianyuan/go-sentry/v2/sentry"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Bool returns a pointer to the bool value.
func Bool(v bool) *bool {
	return &v
}

// Int returns a pointer to the int value.
func Int(v int) *int {
	return &v
}

// checkClientGet returns a `found` bool and an `error` to indicate if a Get request was successful.
// The following return values are meaningful:
// `true`, `nil` => a resource was successfully found
// `false`, `nil` => a resource was successfully not found
// `false`, `err` => encountered an unexpected error
func checkClientGet(resp *sentry.Response, err error, d *schema.ResourceData) (bool, error) {
	if err != nil {
		if resp != nil && resp.Response.StatusCode == http.StatusNotFound {
			d.SetId("")
			return false, nil
		}

		return false, err
	}

	return true, nil
}
