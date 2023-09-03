package godom

import (
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
