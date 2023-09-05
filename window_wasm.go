//go:build wasm

package godom

import (
	"syscall/js"
)

var globalWasmWindow Window

func globalClear() {
	globalWasmWindow = nil
}

func Global() Window {
	if globalWasmWindow == nil {
		globalWasmWindow = &wasmWindow{v: FromJsValue(js.Global())}
	}
	return globalWasmWindow
}

var _ Window = (*wasmWindow)(nil)

type wasmWindow struct {
	node
	v Value
	d *document
	n *navigator
	l *location
	c *console
}

func (w *wasmWindow) Value() Value {
	return w.v
}

func (w *wasmWindow) Document() Document {
	if w.d != nil {
		return w.d
	}
	w.d = &document{}
	w.d.this = w.v.Get("document")
	w.d.this.SetGoValue(w.d)
	return w.d
}

func (w *wasmWindow) Navigator() Navigator {
	if w.n != nil {
		return w.n
	}
	w.n = &navigator{}
	w.n.this = w.v.Get("navigator")
	return w.n
}
func (w *wasmWindow) Location() Location {
	if w.l != nil {
		return w.l
	}
	w.l = &location{}
	w.l.this = w.v.Get("location")
	return w.l
}
func (w *wasmWindow) Console() Console {
	if w.c != nil {
		return w.c
	}
	w.c = &console{}
	w.c.this = w.v.Get("console")
	return w.c
}
