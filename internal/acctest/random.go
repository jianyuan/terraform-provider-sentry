package acctest

import (
	pluginacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func RandomWithPrefix(name string) string {
	return pluginacctest.RandomWithPrefix(name)
}
