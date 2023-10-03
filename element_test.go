//go:build !wasm

package godom

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestElement_ReplaceWith_noWasm(t *testing.T) {

	a := require.New(t)
	api := NewDocApi(Document().This())
	body := api.H(`<body><div id="toReplace"/></body>"`)
	replacement := api.H(`<div id="replacingWith"/>`)
	body.ChildNodes()[0].(Element).ReplaceWith(replacement)
	a.Equal(`<body><div id="replacingWith"/></body>`, body.String())

}
