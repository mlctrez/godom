//go:build wasm

package godom

import (
	"syscall/js"
)

var globalWasmWindow Window

func Global() Window {
	if globalWasmWindow == nil {
		globalWasmWindow = &wasmWindow{v: FromJsValue(js.Global())}
	}
	return globalWasmWindow
}

var _ Window = (*wasmWindow)(nil)

type wasmWindow struct {
	v Value
}

func (w *wasmWindow) Document() Document   { panic(IM) }
func (w *wasmWindow) Navigator() Navigator { panic(IM) }
func (w *wasmWindow) Location() Location   { panic(IM) }
func (w *wasmWindow) Console() Console     { panic(IM) }
