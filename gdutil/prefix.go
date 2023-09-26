package gdutil

import (
	"net/url"
	"path/filepath"
)

func GetPrefix(u *url.URL) string {
	prefix := filepath.Dir(u.Path)
	switch prefix {
	case ".", "/":
		prefix = ""
	}
	return prefix
}
