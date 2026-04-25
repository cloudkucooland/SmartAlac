package sa

import (
	"strings"
)

func mediumFormat(f string) string {
	return strings.ReplaceAll(f, "\"", "″")
}
