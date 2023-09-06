//go:build js && wasm

package godom

import (
	"syscall/js"
)

var _ Value = (*wasmValue)(nil)

func FromJsValue(v js.Value) Value {
	return &wasmValue{jsv: v}
}

func FromJsValues(v ...js.Value) []Value {
	result := make([]Value, len(v))
	for i, value := range v {
		result[i] = FromJsValue(value)
	}
	return result
}

func ToJsValue(arg interface{}) interface{} {
	switch v := arg.(type) {
	case func(this Value, args []Value) interface{}:
		return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return v(FromJsValue(this), FromJsValues(args...))
		})
	case *wasmWindow:
		return ToJsValue(v.v)
	case *wasmValue:
		return v.jsv
	case *navigator:
		return ToJsValue(v.this)
	case *document:
		return ToJsValue(v.this)
	case *location:
		return ToJsValue(v.this)
	case *console:
		return ToJsValue(v.this)
	default:
		return js.ValueOf(v)
	}
}

func ToJsValues(args ...interface{}) (result []interface{}) {
	result = make([]interface{}, len(args))
	for i, arg := range args {
		result[i] = ToJsValue(arg)
	}
	return result
}

type wasmValue struct {
	jsv js.Value
	gov interface{}
}

func (v *wasmValue) GoValue() interface{}      { return v.gov }
func (v *wasmValue) SetGoValue(gv interface{}) { v.gov = gv }

func (v *wasmValue) JSValue() Value          { return v }
func (v *wasmValue) Equal(w Value) bool      { return v.jsv.Equal(w.(*wasmValue).jsv) }
func (v *wasmValue) InstanceOf(t Value) bool { return v.jsv.InstanceOf(t.(*wasmValue).jsv) }
func (v *wasmValue) IsUndefined() bool       { return v.jsv.IsUndefined() }
func (v *wasmValue) IsNull() bool            { return v.jsv.IsNull() }
func (v *wasmValue) IsNaN() bool             { return v.jsv.IsNaN() }
func (v *wasmValue) Type() Type              { return Type(v.jsv.Type()) }
func (v *wasmValue) Length() int             { return v.jsv.Length() }
func (v *wasmValue) Float() float64          { return v.jsv.Float() }
func (v *wasmValue) Int() int                { return v.jsv.Int() }
func (v *wasmValue) Bool() bool              { return v.jsv.Bool() }
func (v *wasmValue) Truthy() bool            { return v.jsv.Truthy() }
func (v *wasmValue) String() string          { return v.jsv.String() }

func (v *wasmValue) Get(p string) Value            { return FromJsValue(v.jsv.Get(p)) }
func (v *wasmValue) Set(p string, x interface{})   { v.jsv.Set(p, ToJsValue(x)) }
func (v *wasmValue) Delete(p string)               { v.jsv.Delete(p) }
func (v *wasmValue) Index(i int) Value             { return FromJsValue(v.jsv.Index(i)) }
func (v *wasmValue) SetIndex(i int, x interface{}) { v.jsv.SetIndex(i, ToJsValue(x)) }
func (v *wasmValue) New(args ...interface{}) Value {
	return FromJsValue(v.jsv.New(ToJsValues(args...)...))
}
func (v *wasmValue) Call(m string, args ...interface{}) Value {
	return FromJsValue(v.jsv.Call(m, ToJsValues(args...)...))
}
func (v *wasmValue) Invoke(args ...interface{}) Value {
	return FromJsValue(v.jsv.Invoke(ToJsValues(args...)...))
}
func (v *wasmValue) Bytes() []byte {
	uint8Array := js.Global().Get("Uint8Array").New(v.jsv)
	dst := make([]byte, uint8Array.Length())
	js.CopyBytesToGo(dst, uint8Array)
	return dst
}
