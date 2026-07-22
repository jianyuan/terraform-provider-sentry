package resourceid

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ImportState1PartIDFunc constructs a single-part state ID from an attribute or Primary.ID.
func ImportState1PartIDFunc(resourceAddress, attributeKey string) func(s *terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		attrMap, err := getResourceAttributes(s, resourceAddress)
		if err != nil {
			return "", err
		}

		val := getAttributeOrID(attrMap, attributeKey)
		if val == "" {
			return "", fmt.Errorf("resource %s attribute %q cannot be empty", resourceAddress, attributeKey)
		}

		return val, nil
	}
}

// ImportState2PartIDFunc constructs a 2-part state ID ("part1/part2") from resource attributes or Primary.ID.
func ImportState2PartIDFunc(resourceAddress, attributeKey1, attributeKey2 string) func(s *terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		attrMap, err := getResourceAttributes(s, resourceAddress)
		if err != nil {
			return "", err
		}

		val1 := getAttributeOrID(attrMap, attributeKey1)
		val2 := getAttributeOrID(attrMap, attributeKey2)

		return BuildPath(val1, val2)
	}
}

// ImportState3PartIDFunc constructs a 3-part state ID ("part1/part2/part3") from resource attributes or Primary.ID.
func ImportState3PartIDFunc(resourceAddress, attributeKey1, attributeKey2, attributeKey3 string) func(s *terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		attrMap, err := getResourceAttributes(s, resourceAddress)
		if err != nil {
			return "", err
		}

		val1 := getAttributeOrID(attrMap, attributeKey1)
		val2 := getAttributeOrID(attrMap, attributeKey2)
		val3 := getAttributeOrID(attrMap, attributeKey3)

		return BuildPath(val1, val2, val3)
	}
}

// ImportStateURL1PartIDFunc constructs a 1-part URL state ID by replacing placeholderA in urlTemplate.
func ImportStateURL1PartIDFunc(resourceAddress, urlTemplate, placeholderA, attributeKeyA string) func(s *terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		attrMap, err := getResourceAttributes(s, resourceAddress)
		if err != nil {
			return "", err
		}

		valA := getAttributeOrID(attrMap, attributeKeyA)
		return Build1(urlTemplate, placeholderA, valA)
	}
}

// ImportStateURL2PartIDFunc constructs a 2-part URL state ID by replacing placeholderA and placeholderB in urlTemplate.
func ImportStateURL2PartIDFunc(resourceAddress, urlTemplate, placeholderA, attributeKeyA, placeholderB, attributeKeyB string) func(s *terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		attrMap, err := getResourceAttributes(s, resourceAddress)
		if err != nil {
			return "", err
		}

		valA := getAttributeOrID(attrMap, attributeKeyA)
		valB := getAttributeOrID(attrMap, attributeKeyB)

		return Build2(urlTemplate, placeholderA, valA, placeholderB, valB)
	}
}

// ImportStateURL3PartIDFunc constructs a 3-part URL state ID by replacing 3 placeholders in urlTemplate.
func ImportStateURL3PartIDFunc(
	resourceAddress, urlTemplate string,
	placeholderA, attributeKeyA string,
	placeholderB, attributeKeyB string,
	placeholderC, attributeKeyC string,
) func(s *terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		attrMap, err := getResourceAttributes(s, resourceAddress)
		if err != nil {
			return "", err
		}

		valA := getAttributeOrID(attrMap, attributeKeyA)
		valB := getAttributeOrID(attrMap, attributeKeyB)
		valC := getAttributeOrID(attrMap, attributeKeyC)

		return Build3(urlTemplate, placeholderA, valA, placeholderB, valB, placeholderC, valC)
	}
}

func getResourceAttributes(s *terraform.State, resourceAddress string) (map[string]string, error) {
	rs, ok := s.RootModule().Resources[resourceAddress]
	if !ok {
		return nil, fmt.Errorf("resource not found in terraform state: %s", resourceAddress)
	}
	if rs.Primary == nil {
		return nil, fmt.Errorf("resource has no primary instance state: %s", resourceAddress)
	}
	return rs.Primary.Attributes, nil
}

// getAttributeOrID retrieves the requested key from the attributes map.
// If key is "id" or "", it falls back to checking the primary ID field.
func getAttributeOrID(attributes map[string]string, key string) string {
	if key == "" || strings.ToLower(key) == "id" {
		return attributes["id"]
	}
	return attributes[key]
}
