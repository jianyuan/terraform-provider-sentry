package acctest

import (
	"regexp"
	"strings"

	"github.com/samber/lo"
)

func ExpectLiteralError(text string) *regexp.Regexp {
	return regexp.MustCompile(strings.Join(lo.Map(strings.Split(text, " "), func(v string, _ int) string {
		return regexp.QuoteMeta(v)
	}), `\s+`))
}
