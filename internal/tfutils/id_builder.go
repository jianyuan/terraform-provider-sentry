package tfutils

import (
	"fmt"
	"strings"
)

func BuildTwoPartId(a, b string) string {
	return fmt.Sprintf("%s/%s", a, b)
}

func SplitTwoPartId(id, a, b string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected %s/%s", id, a, b)
	}
	return parts[0], parts[1], nil
}

func BuildThreePartId(a, b, c string) string {
	return fmt.Sprintf("%s/%s/%s", a, b, c)
}

func SplitThreePartId(id, a, b, c string) (string, string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%s), expected %s/%s/%s", id, a, b, c)
	}
	return parts[0], parts[1], parts[2], nil
}

func BuildFourPartId(a, b, c, d string) string {
	return fmt.Sprintf("%s/%s/%s/%s", a, b, c, d)
}

func SplitFourPartId(id, a, b, c, d string) (string, string, string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 4 || parts[0] == "" || parts[1] == "" || parts[2] == "" || parts[3] == "" {
		return "", "", "", "", fmt.Errorf("unexpected format of ID (%s), expected %s/%s/%s/%s", id, a, b, c, d)
	}
	return parts[0], parts[1], parts[2], parts[3], nil
}
