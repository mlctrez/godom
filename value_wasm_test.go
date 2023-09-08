package godom

import (
	"github.com/stretchr/testify/assert"
	"syscall/js"
	"testing"
)

func TestValue_GetSet(t *testing.T) {
	a := assert.New(t)
	g := Global()

	a.True(g.Get("undefinedValue").IsUndefined())

	g.Set("definedValue", "abc123")
	dv := g.Get("definedValue")
	a.True(!dv.IsUndefined())
	a.True(dv.Type() == TypeString)
	a.Equal("abc123", dv.String())
}

func TestValue_GetSetFunc(t *testing.T) {
	a := assert.New(t)
	g := Global()

	var invokeArgs []Value

	aFunc := ToJsValue(func(this Value, args []Value) interface{} {
		invokeArgs = args
		return nil
	})
	g.Set("testFunction", aFunc)
	dv := g.Get("testFunction")
	a.True(!dv.IsUndefined())
	a.True(dv.Type() == TypeFunction)
	dv.Invoke("testing")

	a.True(len(invokeArgs) == 1)
	a.Equal("testing", invokeArgs[0].String())

	g.Delete("testFunction")
	dv = g.Get("testFunction")
	a.True(dv.IsUndefined())
}

func TestToJsValue(t *testing.T) {
	a := assert.New(t)
	a.IsType(js.Value{}, ToJsValue(Global()))
	a.IsType(js.Value{}, ToJsValue(Document()))
}

func TestWasmValue_GoValue(t *testing.T) {
	a := assert.New(t)
	g := Global()
	g.SetGoValue(g)
	a.IsType(&wasmValue{}, g.GoValue())
}

func TestWasmValue_IsNaN(t *testing.T) {
	a := assert.New(t)
	g := Global()
	g.Set("thisIsNotANumber", g.Get("NaN"))
	a.True(g.Get("thisIsNotANumber").IsNaN())
}

func TestWasmValue_Array(t *testing.T) {
	a := assert.New(t)
	g := Global()
	g.Set("testArray", g.Get("Array").New("a", "b", "c"))
	testArray := g.Get("testArray")
	a.True(testArray.Length() == 3)
	a.True(testArray.Index(0).String() == "a")
	testArray.SetIndex(0, "new")
	a.True(testArray.Index(0).String() == "new")
}

func TestWasmValue_Equal(t *testing.T) {
	a := assert.New(t)
	g := Global()
	g.Set("testArray", g.Get("Array").New("a", "b", "c"))
	testArray := g.Get("testArray")
	a.True(testArray.Equal(g.Get("testArray")))
}

func TestWasmValue_InstanceOf(t *testing.T) {
	a := assert.New(t)
	g := Global()
	g.Set("testArray", g.Get("Array").New("a", "b", "c"))
	testArray := g.Get("testArray")
	a.True(testArray.InstanceOf(g.Get("Array")))
}

func TestWasmValue_Bytes(t *testing.T) {
	a := assert.New(t)
	g := Global()
	g.Set("testArray", g.Get("Array").New(0, 1, 2, 3))
	testArray := g.Get("testArray")
	bytes := testArray.Bytes()
	a.True(len(bytes) == 4)
}

func TestWasmValue_JsValue(t *testing.T) {
	a := assert.New(t)
	g := Global()
	a.True(g.Equal(g.JSValue()))
}

func TestWasmValue_Float(t *testing.T) {
	a := assert.New(t)
	g := Global()
	var floatValue = 123.456
	g.Set("testFloatValue", floatValue)
	a.Equal(floatValue, g.Get("testFloatValue").Float())
}

func TestWasmValue_Int(t *testing.T) {
	a := assert.New(t)
	g := Global()
	var intValue = 123456
	g.Set("testIntValue", intValue)
	a.Equal(intValue, g.Get("testIntValue").Int())
}
