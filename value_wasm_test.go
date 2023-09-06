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
}

func TestValue_ToJsValue(t *testing.T) {
	a := assert.New(t)
	a.IsType(js.Value{}, ToJsValue(Global()))
	a.IsType(js.Value{}, ToJsValue(Global().Location()))
	a.IsType(js.Value{}, ToJsValue(Global().Navigator()))
	a.IsType(js.Value{}, ToJsValue(Global().Document()))
	a.IsType(js.Value{}, ToJsValue(Global().Console()))

}
