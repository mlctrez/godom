package godom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestType_String(t *testing.T) {

	cases := map[Type]string{
		TypeUndefined: "undefined",
		TypeNull:      "null",
		TypeBoolean:   "boolean",
		TypeNumber:    "number",
		TypeString:    "string",
		TypeSymbol:    "symbol",
		TypeObject:    "object",
		TypeFunction:  "function",
		Type(99):      "unknown :99",
	}
	for typ, expected := range cases {
		if typ.String() != expected {
			t.Errorf("expected %q but got %q", expected, typ.String())
		}
	}
}

func TestUtil_ptr(t *testing.T) {

	doc1 := &document{}
	doc2 := &document{}
	if ptr(doc1) != ptr(doc1) {
		t.Error("failed equality test")
	}
	if ptr(doc1) == ptr(doc2) {
		t.Error("failed inequality test")
	}

}

/*
	Everything below is called from both wasm and !wasm test cases
	to ensure that assertions are consistent across architectures.
*/

func testValueNew(t *testing.T) {
	a := assert.New(t)
	obj := Global().Get("Object")
	a.Equal(TypeFunction, obj.Type())
	a.Equal(TypeObject, obj.New().Type())
}

func testValueIsNaN(t *testing.T) {
	a := assert.New(t)
	obj := Global().Get("NaN")
	a.Equal(TypeNumber, obj.Type())
	a.True(obj.IsNaN())
}

func testValueInstanceOf(t *testing.T) {
	a := assert.New(t)
	g := Global()

	p := "theVarUnderTest"
	so := g.Get("Object").New()
	g.Set(p, so)

	fromGlobal := g.Get(p)
	stringType := g.Get("Object")
	a.True(fromGlobal.InstanceOf(stringType), "js.Value.InstanceOf test failed")
}
