package acctest

import (
	"regexp"
	"strings"

	"github.com/jianyuan/go-utils/sliceutils"
)

func ExpectLiteralError(lines ...string) *regexp.Regexp {
	return regexp.MustCompile(strings.Join(sliceutils.Map(func(v string) string {
		return regexp.QuoteMeta(v)
	}, lines), `\n`))
}
