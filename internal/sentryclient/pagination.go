package sentryclient

import (
	"net/http"

	"github.com/jianyuan/go-utils/ptr"
	"github.com/peterhellberg/link"
)

func ParseNextPaginationCursor(resp *http.Response) *string {
	rels := link.ParseResponse(resp)
	if nextRel, ok := rels["next"]; ok && nextRel.Extra["results"] == "true" && nextRel.Extra["cursor"] != "" {
		return ptr.Ptr(nextRel.Extra["cursor"])
	}
	return nil
}
