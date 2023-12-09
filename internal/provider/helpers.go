package provider

import (
	"fmt"
	"strings"
)

func buildTwoPartID(a, b string) string {
	return fmt.Sprintf("%s/%s", a, b)
}

func splitTwoPartID(id, a, b string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected %s/%s", id, a, b)
	}
	return parts[0], parts[1], nil
}

func buildThreePartID(a, b, c string) string {
	return fmt.Sprintf("%s/%s/%s", a, b, c)
}

func splitThreePartID(id, a, b, c string) (string, string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%s), expected %s/%s/%s", id, a, b, c)
	}
	return parts[0], parts[1], parts[2], nil
}

func splitSentryAlertID(id string) (org string, project string, alertID string, err error) {
	org, project, alertID, err = splitThreePartID(id, "organization-slug", "project-slug", "alert-id")
	return
}
