package acctest

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

var _ knownvalue.Check = stringJson{}

type stringJson struct{}

func (v stringJson) CheckValue(other any) error {
	otherVal, ok := other.(string)

	if !ok {
		return fmt.Errorf("expected string value for StringJson check, got: %T", other)
	}

	var otherValJson any
	err := json.Unmarshal([]byte(otherVal), &otherValJson)
	if err != nil {
		return fmt.Errorf("expected JSON value for StringJson check, got: %s, error: %w", otherVal, err)
	}

	return nil
}

func (v stringJson) String() string {
	return "StringJson"
}

func StringJson() knownvalue.Check {
	return stringJson{}
}
