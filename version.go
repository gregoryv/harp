package harp

import (
	_ "embed"
	"strings"
)

// Version returns the last version found in the changelog
func Version() string {
	prefix := "## ["
	from := strings.Index(changelog, prefix) + len(prefix)
	to := from + strings.Index(changelog[from:], "]")
	return changelog[from:to]
}

//go:embed changelog.md
var changelog string
