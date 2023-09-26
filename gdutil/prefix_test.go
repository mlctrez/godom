package gdutil

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestGetPrefix(t *testing.T) {
	a := require.New(t)
	a.Equal("", GetPrefix(&url.URL{Path: ""}))
	a.Equal("", GetPrefix(&url.URL{Path: "/"}))
	a.Equal("", GetPrefix(&url.URL{Path: "/example"}))
	a.Equal("/example", GetPrefix(&url.URL{Path: "/example/"}))
}
