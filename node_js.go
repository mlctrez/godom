//go:build wasm

package godom

import "syscall/js"

func (d *node) AddEventListener(eventType string, fn OnEvent) func() {
	jsFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		fn(FromJsValue(args[0]))
		return nil
	})
	d.this.Call("addEventListener", eventType, jsFunc)
	return func() {
		d.this.Call("removeEventListener", eventType, jsFunc)
		jsFunc.Release()
	}
}
