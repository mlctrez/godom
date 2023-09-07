//go:build js && wasm

package gfet

import (
	"github.com/stretchr/testify/assert"
	"net/textproto"
	"testing"
)

func TestFetch(t *testing.T) {
	a := assert.New(t)

	headers := textproto.MIMEHeader{}
	headers.Set("custom-header", "custom-header-value")

	req := &Request{URL: "/", Headers: headers}
	res, err := Fetch(req)
	a.Nil(err)
	a.True(res.Ok)

}

func TestFetch_rejected(t *testing.T) {
	a := assert.New(t)
	req := &Request{URL: "bad://url"}
	res, err := Fetch(req)
	a.Nil(res)
	a.NotNil(err)
}
