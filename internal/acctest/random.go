package acctest

import (
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func RandomWithPrefix(name string) string {
	return sdkacctest.RandomWithPrefix(name)
}
