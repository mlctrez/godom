//go:build !awasm

package godom

import (
	"bytes"
	"testing"
)

func TestElement_NodeType(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
	e := doc.El("html")
	if e.NodeType() != NodeTypeElement {
		t.Fatalf("expected NodeTypeElement but got")
	}
}

func TestElement_SetAttribute(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
	e := doc.El("html")
	e.SetAttribute("foo", "bar")
	if e.(*element).attributes[0].Name != "foo" {
		t.Fatal("expected attribute not found")
	}
}

func TestElement_Marshal(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
	e := doc.El("html")
	enc := NewEncoder(&bytes.Buffer{})
	e.Marshal(enc)
	if enc.Xml() != "<html/>" {
		t.Fatal("unexpected xml" + enc.Xml())
	}
}

func TestElement_Marshal_Nested(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
	e := doc.El("html")
	e.AppendChild(doc.El("body"))
	enc := NewEncoder(&bytes.Buffer{})
	e.Marshal(enc)
	if enc.Xml() != "<html><body/></html>" {
		t.Fatal("unexpected xml" + enc.Xml())
	}
}

func TestElement_isAlwaysClose(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
	e := doc.El("script")
	if !e.(*element).isAlwaysClose() {
		t.Fatal("isAlwaysClose on script should be true")
	}
}

func TestElement_String(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
	e := doc.El("script")
	if e.String() != "<script></script>" {
		t.Fatal("Element.String test failed")
	}
}

func TestElement_ReplaceWith(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
	e1 := doc.El("p")
	e2 := doc.El("p")
	e1.ReplaceWith(e2)
}

func TestElement_Remove(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
	e1 := doc.El("p")
	e1.Remove()
}

func TestElement_AddEventHandler(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
	e1 := doc.El("p")
	e1.AddEventListener("foo", func(event Value) {})
}
