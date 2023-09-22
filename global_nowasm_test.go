//go:build !wasm

package godom

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestGlobal_Body(t *testing.T) {
	a := require.New(t)
	a.NotPanics(func() { Document().Body() })
	Document().DocumentElement().(*element).children = nil
	a.Panics(func() { Document().Body() })
	Document().This().Call("reset")
}

func TestGlobal_Head(t *testing.T) {
	a := require.New(t)
	a.NotPanics(func() { Document().Head() })
	Document().DocumentElement().(*element).children = nil
	a.Panics(func() { Document().Head() })
	Document().This().Call("reset")
}

func TestGlobal_GetElementById(t *testing.T) {
	a := require.New(t)
	doc := Document()
	h := doc.DocApi().H(`<div id="idOne" match="yes"/>`)
	doc.Body().AppendChild(h)
	a.True(doc.This().Call("getElementById", "foo").IsNull())
	a.False(doc.This().Call("getElementById", "idOne").IsNull())
}

func TestGlobal_AddEventListener(t *testing.T) {

	dummyFunc()

	doc := Document()
	h := doc.DocApi().H(`<div/>`)
	doc.Body().AppendChild(h)
	// just to test that the method can be invoked
	h.AddEventListener("click", func(event Value) {})()
}
