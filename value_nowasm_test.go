//go:build !wasm

package godom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockBasic struct {
	setFunc func(key string, value interface{})
}

func (m *mockBasic) String() string                   { return "String()" }
func (m *mockBasic) LowercaseMethod() string          { return "lowercaseMethod()" }
func (m *mockBasic) MultipleReturn() (string, string) { return "MultipleReturn()", "MultipleReturn()" }
func (m *mockBasic) CheckUndefined()                  {}
func (m *mockBasic) CheckNull() any                   { return nil }

func TestValue_Basic(t *testing.T) {
	a := assert.New(t)
	s := &value{t: TypeObject, v: &mockBasic{}}
	a.Equal(TypeObject, s.Type())
	a.True(s.JSValue().Equal(s))
	func() {
		defer func() { a.Contains(recover(), "method does not exist") }()
		s.Call("MethodNameDoesNotExist")
	}()
	func() {
		defer func() { a.Contains(recover(), "multiple return values") }()
		s.Call("MultipleReturn")
	}()
	a.Equal("String()", s.Call("String").String())
	a.Equal("lowercaseMethod()", s.Call("lowercaseMethod").String())
	a.True(s.Call("CheckUndefined").IsUndefined())
	a.True(s.Call("CheckNull").IsNull())
	func() {
		defer func() { a.Contains(recover(), "not object type") }()
		(&value{t: TypeNull}).Call("notUsed")
	}()
	func() {
		defer func() { a.Contains(recover(), "nil object value") }()
		(&value{t: TypeObject}).Call("notUsed")
	}()

}

type mapObject struct {
	m map[string]interface{}
}

func (mo *mapObject) Get(p string) any    { return toValue(mo.m[p]) }
func (mo *mapObject) Set(p string, x any) { mo.m[p] = x }
func (mo *mapObject) Delete(p string)     { delete(mo.m, p) }

func TestValue_GetSetDelete(t *testing.T) {
	a := assert.New(t)
	mo := &mapObject{m: map[string]interface{}{}}
	s := &value{t: TypeObject, v: mo}
	s.Set("key", "value")
	a.Equal("value", mo.m["key"])
	a.Equal("value", s.Get("key").String())
	s.Delete("key")
	a.Equal(nil, mo.m["key"])
}

type indexObject struct {
	array []any
}

func (mo *indexObject) Index(i int) interface{}       { return mo.array[i] }
func (mo *indexObject) SetIndex(i int, x interface{}) { mo.array[i] = x }
func (mo *indexObject) Length() int                   { return len(mo.array) }

func TestValue_IndexObject(t *testing.T) {
	a := assert.New(t)
	mo := &indexObject{array: []any{"a", "b", "c"}}
	s := toValue(mo)
	a.Equal("a", s.Index(0).String())
	a.Equal("b", s.Index(1).String())
	a.Equal(3, s.Length())
	s.SetIndex(1, "b_new")
	a.Equal("b_new", s.Index(1).String())
}

func TestValue_Numeric(t *testing.T) {
	a := assert.New(t)
	s := toValue(10)
	a.Equal(TypeNumber, s.Type())
	a.Equal(10, s.Int())
	a.Equal(float64(10), s.Float())

	s = toValue(10.1234)
	a.Equal(TypeNumber, s.Type())
	a.Equal(10, s.Int())
	a.Equal(10.1234, s.Float())
}

func TestValue_Bool(t *testing.T) {
	a := assert.New(t)
	a.Equal(TypeBoolean, toValue(true).Type())
	a.Equal(true, toValue(true).Bool())
	a.Equal(false, toValue(false).Bool())
	func() {
		defer func() { a.Contains(recover(), "not type boolean") }()
		(&value{t: TypeObject}).Bool()
	}()
}

/*
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
*/
func TestValue_Truthy(t *testing.T) {

	a := assert.New(t)
	a.False((&value{t: TypeUndefined}).Truthy())
	a.False((&value{t: TypeNull}).Truthy())
	a.False((&value{t: TypeBoolean, v: false}).Truthy())
	a.True((&value{t: TypeBoolean, v: true}).Truthy())
	a.True((&value{t: TypeNumber, v: 1}).Truthy())
	a.True((&value{t: TypeObject, v: "dummy"}).Truthy())

}
