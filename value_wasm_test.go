package godom

import (
	"github.com/stretchr/testify/assert"
	"syscall/js"
	"testing"
)

func TestValue_GetSet(t *testing.T) {
	a := assert.New(t)
	global := Global().Value()

	a.True(global.Get("undefinedValue").IsUndefined())

	global.Set("definedValue", "abc123")
	dv := global.Get("definedValue")
	a.True(!dv.IsUndefined())
	a.True(dv.Type() == TypeString)
	a.Equal("abc123", dv.String())
}

func TestValue_GetSetFunc(t *testing.T) {
	a := assert.New(t)
	global := Global().Value()

	var invokeArgs []Value

	aFunc := ToJsValue(func(this Value, args []Value) interface{} {
		invokeArgs = args
		return nil
	})
	global.Set("testFunction", aFunc)
	dv := global.Get("testFunction")
	a.True(!dv.IsUndefined())
	a.True(dv.Type() == TypeFunction)
	dv.Invoke("testing")

	a.True(len(invokeArgs) == 1)
	a.Equal("testing", invokeArgs[0].String())

	global.Delete("testFunction")
	dv = global.Get("testFunction")
	a.True(dv.IsUndefined())
}

func TestToJsValue(t *testing.T) {
	a := assert.New(t)
	g := Global()
	a.IsType(js.Value{}, ToJsValue(g))
	a.IsType(js.Value{}, ToJsValue(g.Location()))
	a.IsType(js.Value{}, ToJsValue(g.Navigator()))
	a.IsType(js.Value{}, ToJsValue(g.Document()))
	a.IsType(js.Value{}, ToJsValue(g.Console()))
}

func TestWasmValue_GoValue(t *testing.T) {
	a := assert.New(t)
	g := Global()
	value := g.Value()
	value.SetGoValue(g)
	a.IsType(&wasmWindow{}, value.GoValue())
}

func TestWasmValue_IsNaN(t *testing.T) {
	a := assert.New(t)
	value := Global().Value()
	value.Set("thisIsNotANumber", value.Get("NaN"))
	a.True(value.Get("thisIsNotANumber").IsNaN())
}

func TestWasmValue_Array(t *testing.T) {
	a := assert.New(t)
	value := Global().Value()
	value.Set("testArray", value.Get("Array").New("a", "b", "c"))
	testArray := value.Get("testArray")
	a.True(testArray.Length() == 3)
	a.True(testArray.Index(0).String() == "a")
	testArray.SetIndex(0, "new")
	a.True(testArray.Index(0).String() == "new")
}

func TestWasmValue_Equal(t *testing.T) {
	a := assert.New(t)
	value := Global().Value()
	value.Set("testArray", value.Get("Array").New("a", "b", "c"))
	testArray := value.Get("testArray")
	a.True(testArray.Equal(value.Get("testArray")))
}

func TestWasmValue_InstanceOf(t *testing.T) {
	a := assert.New(t)
	value := Global().Value()
	value.Set("testArray", value.Get("Array").New("a", "b", "c"))
	testArray := value.Get("testArray")
	a.True(testArray.InstanceOf(value.Get("Array")))
}

func TestWasmValue_Bytes(t *testing.T) {
	a := assert.New(t)
	value := Global().Value()
	value.Set("testArray", value.Get("Array").New(0, 1, 2, 3))
	testArray := value.Get("testArray")
	bytes := testArray.Bytes()
	a.True(len(bytes) == 4)
}

func TestWasmValue_JsValue(t *testing.T) {
	a := assert.New(t)
	value := Global().Value()
	a.True(value.Equal(value.JSValue()))
}

func TestWasmValue_Float(t *testing.T) {
	a := assert.New(t)
	g := Global().Value()
	var floatValue = 123.456
	g.Set("testFloatValue", floatValue)
	a.Equal(floatValue, g.Get("testFloatValue").Float())
}

func TestWasmValue_Int(t *testing.T) {
	a := assert.New(t)
	g := Global().Value()
	var intValue = 123456
	g.Set("testIntValue", intValue)
	a.Equal(intValue, g.Get("testIntValue").Int())
}
