package acctest

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jianyuan/terraform-provider-sentry/internal/tfutils"
)

func TwoPartImportStateIdFunc(resourceAddress string, part1 string) func(s *terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceAddress]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceAddress)
		}
		return tfutils.BuildTwoPartId(
			rs.Primary.Attributes[part1],
			rs.Primary.ID,
		), nil
	}
}

func ThreePartImportStateIdFunc(resourceAddress string, part1 string, part2 string) func(s *terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceAddress]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceAddress)
		}
		return tfutils.BuildThreePartId(
			rs.Primary.Attributes[part1],
			rs.Primary.Attributes[part2],
			rs.Primary.ID,
		), nil
	}

}

func FourPartImportStateIdFunc(resourceAddress string, part1 string, part2 string, part3 string) func(s *terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceAddress]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceAddress)
		}
		return tfutils.BuildFourPartId(
			rs.Primary.Attributes[part1],
			rs.Primary.Attributes[part2],
			rs.Primary.Attributes[part3],
			rs.Primary.ID,
		), nil
	}
}
