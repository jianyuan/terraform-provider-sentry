package acctest

import (
	"regexp"
	"strings"

	"github.com/samber/lo"
)

func ExpectLiteralError(lines ...string) *regexp.Regexp {
	return regexp.MustCompile(strings.Join(lo.Map(lines, func(v string, _ int) string {
		return regexp.QuoteMeta(v)
	}), `\n`))
}
