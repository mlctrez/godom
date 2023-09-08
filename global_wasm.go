//go:build wasm

package godom

import "syscall/js"

func global() Value {
	return &wasmValue{jsv: js.Global()}
}
