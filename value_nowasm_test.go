//go:build !wasm

package godom

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValue_Null(t *testing.T) {
	v := ToValue(nil)
	if v.Type() != TypeNull {
		t.Fatal("null type incorrect")
	}
	if !v.IsNull() {
		t.Fatal("IsNull() incorrect")
	}
}

func TestValue_Bool(t *testing.T) {
	b := true
	v := ToValue(b)
	if v.Type() != TypeBoolean {
		t.Fatal("bool type incorrect")
	}
	if v.Bool() != true {
		t.Fatal("bool value incorrect")
	}
}

func TestValue_Number(t *testing.T) {
	b := 12345
	v := ToValue(b)
	if v.Type() != TypeNumber {
		t.Fatal("number type incorrect")
	}
	if v.Int() != b {
		t.Fatal("number value incorrect")
	}
}

func TestValue_String(t *testing.T) {
	a := assert.New(t)
	b := "a string"
	v := ToValue(b)
	if v.Type() != TypeString {
		t.Fatal("string type incorrect")
	}
	if v.String() != "a string" {
		t.Fatal("string value incorrect")
	}
	a.Equal("undefined", (&value{t: TypeUndefined}).String())
	a.Equal("null", (&value{t: TypeNull}).String())
	a.Equal("map[]", (&value{t: TypeObject}).String())

	func() {
		defer func() { a.NotNil(recover()) }()
		(&value{t: TypeFunction}).String()
	}()
}

func TestValue_Func(t *testing.T) {
	b := func() {}
	v := ToValue(b)
	if v.Type() != TypeFunction {
		t.Fatal("func type incorrect")
	}
	b2 := func(i int) {}
	v = ToValue(b2)
	if v.Type() != TypeFunction {
		t.Fatal("func type incorrect")
	}
}

func TestValue_Value(t *testing.T) {
	b := true
	v1 := ToValue(b)
	v2 := ToValue(v1)
	if !v1.Equal(v2) {
		t.Fatal("value from value failed")
	}
}

func TestValue_JSValue(t *testing.T) {
	b := true
	v1 := ToValue(b)
	v2 := v1.JSValue()
	if !v1.Equal(v2) {
		t.Fatal("value from value failed")
	}
}

func TestValue_Object(t *testing.T) {
	v := ToValue(Attribute{Name: "name", Value: "value"})
	if v.Type() != TypeObject {
		t.Fatal("incorrect type")
	}
}

func TestValue_Array(t *testing.T) {

	ary := []interface{}{"a", "b", "b"}
	v := ToValue(ary)
	if v.Type() != TypeObject {
		t.Fatal("type incorrect")
	}
	if v.Length() != 3 {
		t.Fatal("length incorrect")
	}
	if v.Index(0).String() != "a" {
		t.Fatal("Index(0) failed")
	}

	v.SetIndex(0, "aa")
	if v.Index(0).String() != "aa" {
		t.Fatal("Index(0) failed")
	}
}

func TestValue_LengthNegativePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("length did not panic with -1")
		}
	}()
	ary := []interface{}{"a", "b", "b"}
	ToValue(ary).SetIndex(-1, "a")
}

func TestValue_LengthOutOfRange(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("length did not panic on out of range")
		}
	}()
	ary := []interface{}{"a", "b", "b"}
	ToValue(ary).SetIndex(10, "a")
}

func TestValue_IndexNegativePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("length did not panic with -1")
		}
	}()
	ary := []interface{}{"a", "b", "b"}
	ToValue(ary).Index(-1)
}

func TestValue_IndexOutOfRange(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("length did not panic on out of range")
		}
	}()
	ary := []interface{}{"a", "b", "b"}
	ToValue(ary).Index(10)
}

func TestValue_LengthNotObject(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("no panic for non object length")
		}
	}()
	ToValue("A").Length()
}

func TestValue_LengthNotArray(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("no panic for value is not an array")
		}
	}()
	(&value{t: TypeObject}).Length()
}

func TestValue_Call(t *testing.T) {
	a := assert.New(t)

	var argPassed string
	v := ToValue(map[string]interface{}{"targetFunction": func(s string) {
		argPassed = s
	}})

	v.Call("targetFunction", "theValue")
	a.Equal("theValue", argPassed)
}

func TestValue_Call_panics(t *testing.T) {
	a := assert.New(t)
	func() {
		defer func() { a.NotNil(recover()) }()
		v := &value{t: TypeString}
		v.Call("notObject", "")
	}()

	func() {
		defer func() { a.NotNil(recover()) }()
		v := &value{t: TypeObject}
		v.Call("missingFunction", "")
	}()
}

func TestValue_NotImplemented_panics(t *testing.T) {
	a := assert.New(t)
	nv := func() Value { return &value{t: TypeObject} }
	func() {
		defer func() { a.NotNil(recover()) }()
		nv().Invoke()
	}()

	func() {
		defer func() { a.NotNil(recover()) }()
		nv().New()
	}()

	func() {
		defer func() { a.NotNil(recover()) }()
		nv().IsUndefined()
	}()

	func() {
		defer func() { a.NotNil(recover()) }()
		nv().IsNaN()
	}()

	func() {
		defer func() { a.NotNil(recover()) }()
		nv().Delete("")
	}()

	func() {
		defer func() { a.NotNil(recover()) }()
		nv().InstanceOf(nil)
	}()

	func() {
		defer func() { a.NotNil(recover()) }()
		nv().Bytes()
	}()

	func() {
		defer func() { a.NotNil(recover()) }()
		nv().Int()
	}()

}

func TestValue_Truthy(t *testing.T) {
	a := assert.New(t)
	a.False((&value{t: TypeUndefined}).Truthy())
	a.False((&value{t: TypeNull}).Truthy())
	a.True(ToValue(true).Truthy())
	a.True(ToValue(1).Truthy())
}

func TestValue_Float(t *testing.T) {
	a := assert.New(t)
	a.Equal(1.234, ToValue(1.234).Float())

	func() {
		defer func() { a.NotNil(recover()) }()
		ToValue("not a float").Float()
	}()
}

func TestInvoke(t *testing.T) {
	buff := &bytes.Buffer{}
	Invoke(buff, "WriteString", "testing")
	assert.Equal(t, "testing", buff.String())
}

func TestValue_set_nil(t *testing.T) {
	a := assert.New(t)
	v := &value{}
	// should not fail for nil v.data
	v.Set("a", "b")
	a.Equal("b", v.data["a"])
}

func TestToValue_Value_slice(t *testing.T) {
	a := assert.New(t)
	toValue := ToValue([]Value{})
	a.IsType([]Value{}, toValue.(*value).data["array"])
}

func TestToValue_panic(t *testing.T) {
	a := assert.New(t)
	defer func() { a.NotNil(recover()) }()
	ToValue([]float64{})
}

func TestValue_GoValue(t *testing.T) {
	a := assert.New(t)
	v := valueT(TypeObject)
	v.SetGoValue(v)
	a.IsType(v, v.GoValue())
}
