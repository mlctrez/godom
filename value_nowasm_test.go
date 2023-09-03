//go:build !wasm

package godom

import "testing"

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
	b := "a string"
	v := ToValue(b)
	if v.Type() != TypeString {
		t.Fatal("string type incorrect")
	}
	if v.String() != "a string" {
		t.Fatal("string value incorrect")
	}
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
