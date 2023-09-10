//go:build !wasm

package godom

var _ Func = (*noWasmFunc)(nil)

type noWasmFunc struct {
	Value
	fn FuncSignature
}

func (nw *noWasmFunc) Release() {
}

func funcOf(fn FuncSignature) Func {
	nwf := &noWasmFunc{fn: fn}
	nwf.Value = toValue(nwf)
	return nwf
}
