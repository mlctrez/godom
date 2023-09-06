//go:build !js && !wasm

package gws

import (
	"github.com/mlctrez/godom"
	"testing"
)

func TestNoWasm_addEventListener(t *testing.T) {
	w := &webSocket{}
	release := w.addEventListener("dummy", func(event godom.Value) {})
	release()
}

func TestNoWasm_sendBinary(t *testing.T) {
	w := &webSocket{}
	_ = w.sendBinary(nil)
}

func TestNoWasm_sendText(t *testing.T) {
	w := &webSocket{}
	_ = w.sendText("")
}
