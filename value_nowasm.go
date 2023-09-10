//go:build !wasm

package godom

import (
	"fmt"
	"github.com/mlctrez/godom/convert"
	"reflect"
	"unicode"
)

var _ Value = (*value)(nil)

type value struct {
	t   Type
	v   any
	gov any
}

func (m *value) GoValue() interface{}      { return m.gov }
func (m *value) SetGoValue(gv interface{}) { m.gov = gv }

// reflectTypeOf returns reflect.TypeOf(in) or nil for nil in.
func reflectTypeOf(in any) (out any) {
	if in != nil {
		out = reflect.TypeOf(in)
	}
	return out
}

func (m *value) Type() Type {
	return m.t
}

func capitalize(str string) string {
	if str == "" {
		return str
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func (m *value) Call(method string, args ...interface{}) Value {
	validFuncTarget(m)
	if method == "" {
		panic(fmt.Sprintf("call type=%q value=%v method=%q method cannot be empty", m.t, reflectTypeOf(m.v), method))
	}
	upMethod := capitalize(method)

	reflectMethod := reflect.ValueOf(m.v).MethodByName(upMethod)
	if !reflectMethod.IsValid() {
		if gt, ok := m.v.(*globalThis); ok {
			if val, ok := gt.Call(method, args); ok {
				return val
			}
		}
		panic(fmt.Sprintf("call type=%q value=%v method=%q method does not exist", m.t, reflectTypeOf(m.v), upMethod))
	}

	var reflectValues []reflect.Value
	if args == nil {
		reflectValues = reflectMethod.Call(nil)
	} else {
		reflectValues = reflectMethod.Call(convert.ToReflectArgs(args...))
	}

	returnValues := convert.FromReflectArgs(reflectValues)
	switch len(returnValues) {
	case 0:
		return &value{t: TypeUndefined, v: nil}
	case 1:
		return toValue(returnValues[0])
	default:
		panic(fmt.Sprintf("call type=%q value=%v method=%q multiple return values", m.t, reflectTypeOf(m.v), upMethod))
	}
}

func validFuncTarget(m *value) {
	switch m.t {
	case TypeFunction, TypeObject:
	default:
		panic(fmt.Sprintf("call type=%q value=%+v not object type", m.t, m.v))
	}
	if m.v == nil {
		panic(fmt.Sprintf("call type=%q value=%+v nil object value", m.t, m.v))
	}
}

func toValues(values []any) []Value {
	result := make([]Value, len(values))
	for i, a := range values {
		result[i] = toValue(a)
	}
	return result
}

func toValue(in any) Value {
	switch vt := in.(type) {
	case *value:
		return vt
	case nil:
		return &value{t: TypeNull, v: nil}
	case string:
		return &value{t: TypeString, v: vt}
	case int, int8, int32, int64, float32, float64:
		return &value{t: TypeNumber, v: convert.ToFloat(vt)}
	case bool:
		return &value{t: TypeBoolean, v: vt}
	case Func:
		return &value{t: TypeFunction, v: vt}
	default:
		return &value{t: TypeObject, v: vt}
	}
}

func (m *value) Invoke(args ...interface{}) Value {
	if m.t != TypeFunction {
		panic("invoke on non-function type")
	}
	if nwf, ok := m.v.(*noWasmFunc); ok {
		return toValue(nwf.fn(nwf, toValues(args)))
	}
	panic(fmt.Sprintf("invoke type=%q value=%+v incorrect value", m.t, m.v))
}

func (m *value) New(args ...interface{}) Value {
	if m.t != TypeFunction {
		panic(fmt.Sprintf("new called on wrong type %q", m.t))
	}
	if args == nil {
		return m.Call("New")
	} else {
		return m.Call("New", args...)
	}
}

func (m *value) Equal(other Value) bool { return ptr(other) == ptr(m) }
func (m *value) JSValue() Value         { return m }

func (m *value) IsUndefined() bool { return m.t == TypeUndefined }
func (m *value) IsNull() bool      { return m.t == TypeNull }

type ValueNaN struct{}

func (m *value) IsNaN() bool {
	return reflectTypeOf(m.v) == reflectTypeOf(ValueNaN{})
}

func (m *value) Get(p string) Value            { return m.Call("Get", p) }
func (m *value) Set(p string, x interface{})   { m.Call("Set", p, x) }
func (m *value) Delete(p string)               { m.Call("Delete", p) }
func (m *value) Index(i int) Value             { return m.Call("Index", i) }
func (m *value) SetIndex(i int, x interface{}) { m.Call("SetIndex", i, x) }
func (m *value) Length() int                   { return m.Call("Length").Int() }
func (m *value) Float() float64                { return convert.ToFloat(m.v) }
func (m *value) Int() int                      { return convert.ToInt(m.v) }
func (m *value) Truthy() bool {
	switch m.t {
	case TypeUndefined, TypeNull:
		return false
	case TypeBoolean:
		return m.Bool()
	case TypeNumber:
		return m.Int() > 0
	default:
		return true
	}
}
func (m *value) Bool() bool {
	if m.t == TypeBoolean {
		return m.v.(bool)
	}
	panic(fmt.Sprintf("call type=%q value=%v not type boolean", m.t, reflectTypeOf(m.v)))
}

func (m *value) String() string { return fmt.Sprintf("%v", m.v) }
func (m *value) Bytes() []byte  { panic(IM) }
func (m *value) InstanceOf(t Value) bool {
	return t.Call("IsInstance", m).Bool()
}
