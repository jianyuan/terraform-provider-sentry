package acctest

import (
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

var _ knownvalue.Check = stringConformingJsonSchema{}

type stringConformingJsonSchema struct {
	jsonSchema *jsonschema.Resolved
}

func (v stringConformingJsonSchema) CheckValue(other any) error {
	otherVal, ok := other.(string)

	if !ok {
		return fmt.Errorf("expected string value for StringConformingJsonSchema check, got: %T", other)
	}

	var otherValJson any
	err := json.Unmarshal([]byte(otherVal), &otherValJson)
	if err != nil {
		return fmt.Errorf("expected JSON value for StringConformingJsonSchema check, got: %s, error: %w", otherVal, err)
	}

	return v.jsonSchema.Validate(otherValJson)
}

func (v stringConformingJsonSchema) String() string {
	return v.jsonSchema.Schema().String()
}

func StringConformingJsonSchema(jsonSchema *jsonschema.Resolved) knownvalue.Check {
	return stringConformingJsonSchema{jsonSchema: jsonSchema}
}
