package godom

import "syscall/js"

var _ Func = (*wasmFunc)(nil)

type wasmFunc struct {
	jsf js.Func
}

func (w wasmFunc) Release() {
	w.jsf.Release()
}

func funcOf(fn func(this Value, args []Value) any) Func {
	jsf := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return fn(FromJsValue(this), FromJsValues(args...))
	})
	return &wasmFunc{jsf: jsf}
}
