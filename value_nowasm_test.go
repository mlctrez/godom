//go:build !wasm

package godom

import (
	"fmt"
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

	a.Contains(recoverString(func() { s.Call("NotExist") }), "method does not exist")
	a.Contains(recoverString(func() { s.Call("MultipleReturn") }), "multiple return values")

	a.Equal("String()", s.Call("String").String())
	a.Equal("lowercaseMethod()", s.Call("lowercaseMethod").String())
	a.True(s.Call("CheckUndefined").IsUndefined())
	a.True(s.Call("CheckNull").IsNull())

	v := &value{t: TypeNull}
	a.Contains(recoverString(func() { v.Call("method") }), "not object type")
	v = &value{t: TypeObject}
	a.Contains(recoverString(func() { v.Call("method") }), "nil object value")
	v = &value{t: TypeObject, v: ""}
	a.Contains(recoverString(func() { v.Call("") }), "method cannot be empty")

}

func TestValue_GetSetDelete(t *testing.T) {
	a := assert.New(t)
	mo := &ValueObject{m: map[string]interface{}{}}
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

func TestValue_Truthy(t *testing.T) {
	// TODO: validate against js, this is a best guess for now
	a := assert.New(t)
	a.False((&value{t: TypeUndefined}).Truthy())
	a.False((&value{t: TypeNull}).Truthy())
	a.False((&value{t: TypeBoolean, v: false}).Truthy())
	a.True((&value{t: TypeBoolean, v: true}).Truthy())
	a.True((&value{t: TypeNumber, v: 1}).Truthy())
	a.True((&value{t: TypeObject, v: "dummy"}).Truthy())
}

func TestValue_funcChecks(t *testing.T) {
	a := assert.New(t)
	s := &value{t: TypeObject}
	a.Contains(recoverString(func() { s.Invoke() }), "invoke on non-function type")

	s = &value{t: TypeFunction}
	a.Contains(recoverString(func() { s.Invoke() }), "incorrect value")
}

func recoverString(toTest func()) (result string) {
	defer func() { result = fmt.Sprintf("%v", recover()) }()
	toTest()
	return ""
}

func TestValue_GoValue(t *testing.T) {
	a := assert.New(t)
	g := Global()
	g.SetGoValue(g)
	a.IsType(&value{}, g.GoValue())

}

type newWithArgs struct {
	ValueObject
	args []any
}

func (n *newWithArgs) New(args ...any) any {
	return &newWithArgs{args: args}
}

func (n *newWithArgs) Args() any {
	return toValue(n.args)
}

func TestValue_New(t *testing.T) {
	testValueNew(t)
	// additional testing for new with multiple args
	g := Global()
	g.Set("NewWithArgs", &value{t: TypeFunction, v: &newWithArgs{}})
	v := g.Get("NewWithArgs").New("a", "b")
	a := assert.New(t)
	a.IsType(&newWithArgs{}, v.(*value).v)
	call := v.Call("Args")
	// TODO: revisit this after creating ValueArray
	a.EqualValues([]any{"a", "b"}, call.(*value).v)
}

func TestValue_New_error(t *testing.T) {
	a := assert.New(t)
	a.Contains(recoverString(func() { (&value{t: TypeObject}).New() }), "new called on wrong type")
}

func Test_capitalize(t *testing.T) {
	// should not panic on runes[0]
	capitalize("")
	a := assert.New(t)
	a.Equal("Test", capitalize("test"))
}

func TestValue_IsNaN(t *testing.T) {
	testValueIsNaN(t)
}

func TestValue_InstanceOf(t *testing.T) {
	testValueInstanceOf(t)
}

func TestValue_Bytes(t *testing.T) {
	a := assert.New(t)
	a.Contains(recoverString(func() {
		(&value{}).Bytes()
	}), "not implemented")
}
