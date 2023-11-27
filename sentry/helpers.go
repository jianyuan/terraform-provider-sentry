package sentry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/deste-org/terraform-provider-sentry/sentry/lib"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func buildTwoPartID(a, b string) string {
	return fmt.Sprintf("%s/%s", a, b)
}

func splitTwoPartID(id, a, b string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected %s/%s", id, a, b)
	}
	return parts[0], parts[1], nil
}

func buildThreePartID(a, b, c string) string {
	return fmt.Sprintf("%s/%s/%s", a, b, c)
}

func splitThreePartID(id, a, b, c string) (string, string, string, error) {
	parts := strings.SplitN(id, "/", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%s), expected %s/%s/%s", id, a, b, c)
	}
	return parts[0], parts[1], parts[2], nil
}

func splitSentryAlertID(id string) (org string, project string, alertID string, err error) {
	org, project, alertID, err = splitThreePartID(id, "organization-slug", "project-slug", "alert-id")
	return
}

func SuppressEquivalentJSONDiffs(k, old, new string, d *schema.ResourceData) bool {
	var o interface{}
	if err := json.Unmarshal([]byte(old), &o); err != nil {
		return false
	}

	var n interface{}
	if err := json.Unmarshal([]byte(new), &n); err != nil {
		return false
	}

	return reflect.DeepEqual(o, n)
}

// followShape reshapes the value into the provided shape
func followShape(shape, value interface{}) interface{} {
	switch shape := shape.(type) {
	case map[string]interface{}:
		value, ok := interface{}(value).(map[string]interface{})
		if !ok {
			return nil
		}

		v := make(map[string]interface{})
		for k, shapeValue := range shape {
			v[k] = followShape(shapeValue, value[k])
		}
		return v
	case []interface{}:
		value, ok := interface{}(value).([]interface{})
		if !ok {
			return nil
		}

		v := make([]interface{}, 0, len(shape))
		for i, shapeValue := range shape {
			v = append(v, followShape(shapeValue, value[i]))
		}
		return v
	default:
		return value
	}
}

func flattenStringSet(strings []string) *schema.Set {
	flattenedStrings := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range strings {
		flattenedStrings.Add(v)
	}
	return flattenedStrings
}

func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		if val, ok := v.(string); ok && val != "" {
			vs = append(vs, val)
		}
	}
	return vs
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
