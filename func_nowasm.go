//go:build !wasm

package godom

var _ Func = (*noWasmFunc)(nil)

type noWasmFunc struct {
}

func (nw *noWasmFunc) Release() {

}

func funcOf(func(this Value, args []Value) any) Func {
	return &noWasmFunc{}
}
