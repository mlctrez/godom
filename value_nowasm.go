//go:build !wasm

package godom

import (
	"fmt"
	"reflect"

	"github.com/mlctrez/godom/convert"
)

var _ Value = (*value)(nil)

type value struct {
	t Type
	// data contains the backing value
	data map[string]interface{}
	gov  interface{}
}

func (v *value) GoValue() interface{}      { return v.gov }
func (v *value) SetGoValue(gv interface{}) { v.gov = gv }

func valueT(t Type) *value { return &value{t: t} }

func (v *value) set(key string, val interface{}) *value {
	if v.data == nil {
		v.data = map[string]interface{}{key: val}
	} else {
		v.data[key] = val
	}
	return v
}

func ToValue(i interface{}) Value {

	switch v := i.(type) {
	case Value:
		return v
	case *element:
		return v.this
	case nil:
		return valueT(TypeNull)
	case bool:
		return valueT(TypeBoolean).set("bool", v)
	case int, int8, int16, int32, int64, float32, float64:
		return valueT(TypeNumber).set("number", v)
	case string:
		return valueT(TypeString).set("string", v)
	case map[string]interface{}:
		return &value{t: TypeObject, data: v}
	case []Value:
		return valueT(TypeObject).set("array", v)
	case []interface{}:
		return valueT(TypeObject).set("array", v)
	case []reflect.Value:
		args := convert.FromReflectArgs(v)
		if len(args) == 1 {
			return ToValue(args[0])
		}
		return valueT(TypeObject).set("array", args)
	}

	if reflect.ValueOf(i).Kind() == reflect.Func {
		return valueT(TypeFunction).set("func", i)
	}
	if reflect.ValueOf(i).Kind() == reflect.Struct {
		return valueT(TypeObject).set("object", i)
	}
	panic(fmt.Errorf("ToValue conversion failed for kind %s %s", reflect.TypeOf(i).Kind(), i))
}

func (v *value) JSValue() Value              { return v }
func (v *value) Equal(w Value) bool          { return ptr(w) == ptr(v) }
func (v *value) IsUndefined() bool           { panic(IM) }
func (v *value) IsNull() bool                { return v.t == TypeNull }
func (v *value) IsNaN() bool                 { panic(IM) }
func (v *value) Type() Type                  { return v.t }
func (v *value) Get(p string) Value          { return ToValue(v.data[p]) }
func (v *value) Set(p string, x interface{}) { v.set(p, x) }
func (v *value) Delete(p string)             { panic(IM) }

func (v *value) Index(i int) Value {
	if i < 0 || i+1 > v.Length() {
		panic("index out of bounds")
	}
	return ToValue(v.data["array"].([]interface{})[i])
}

func (v *value) SetIndex(i int, x interface{}) {
	if i < 0 || i+1 > v.Length() {
		panic("index out of bounds")
	}
	v.data["array"].([]interface{})[i] = x
}

func (v *value) Length() int {
	if v.Type() != TypeObject {
		panic("value is not an object")
	}
	if ary, ok := v.data["array"].([]interface{}); ok {
		return len(ary)
	}
	panic("value is not an array")
}

func (v *value) Call(m string, args ...interface{}) Value {
	if v.Type() != TypeObject {
		panic("value is not object")
	}
	if f, ok := v.data[m]; ok {
		of := reflect.ValueOf(f)
		reflectArgs := convert.ToReflectArgs(args...)
		call := of.Call(reflectArgs)
		return ToValue(call)
	} else {
		panic(fmt.Errorf("no such method %q", m))
	}
}

func (v *value) Invoke(args ...interface{}) Value { panic(IM) }
func (v *value) New(args ...interface{}) Value    { panic(IM) }

func (v *value) Float() float64 {
	switch v.t {
	case TypeNumber:
		return convert.ToFloat(v.data["number"])
	}
	panic(fmt.Errorf("value %s converstion to Float() failed", v.t))
}

func (v *value) Int() int {
	switch v.t {
	case TypeNumber:
		return convert.ToInt(v.data["number"])
	}
	panic(fmt.Errorf("value %s converstion to Int() failed", v.t))
}

func (v *value) Bool() bool {
	return v.data["bool"].(bool)
}

func (v *value) Truthy() bool {
	switch v.t {
	case TypeUndefined, TypeNull:
		return false
	case TypeBoolean:
		return v.Bool()
	case TypeNumber:
		return v.Int() > 0
	default:
		return true
	}
}

func (v *value) String() string {
	switch v.t {
	case TypeUndefined, TypeNull:
		return v.t.String()
	case TypeString:
		return v.data["string"].(string)
	case TypeObject:
		return fmt.Sprintf("%+v", v.data)
	}
	panic(fmt.Errorf("type %d String() not implemented", v.t))
}

func (v *value) InstanceOf(t Value) bool { panic(IM) }

func (v *value) Bytes() []byte { panic(IM) }

// Invoke is a sample demonstrating go reflection.
func Invoke(any interface{}, name string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)
}
