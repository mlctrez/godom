//go:build !wasm

package godom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGlobal_Object(t *testing.T) {
	a := assert.New(t)
	get := Global().Get("Object")
	a.IsType(&value{}, get)
	a.Equal(TypeFunction, get.Type())
	a.IsType(&ValueObject{}, get.(*value).v)

	// Object should not allow get set
	a.Contains(recoverString(func() {
		get.Set("test", "bad")
	}), "assignment to entry in nil map")
}

func TestMockObject_GetAttribute(t *testing.T) {
	a := assert.New(t)
	obj := &mockObject{}
	obj.SetAttribute("testAttribute", "value")
	a.Equal("value", obj.GetAttribute("testAttribute").String())
}

func TestMockObject_Call_notFound(t *testing.T) {
	a := assert.New(t)
	obj := &globalThis{}
	call, b := obj.Call("notThere", nil)
	a.Nil(call)
	a.False(b)
}
