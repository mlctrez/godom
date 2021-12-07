package godom

import (
	"bytes"
	"testing"
)

func TestElement_NodeType(t *testing.T) {
	e := &element{}
	if e.NodeType() != NodeTypeElement {
		t.Fatalf("expected NodeTypeElement but got")
	}
}

func TestElement_SetAttribute(t *testing.T) {
	e := &element{}
	e.SetAttribute("foo", "bar")
	if e.attributes[0].Name != "foo" {
		t.Fatal("expected attribute not found")
	}
}

func TestElement_Remove(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Element.Remove did not panic")
		}
	}()
	e := &element{}
	e.Remove()
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
	e := &element{}
	e.nodeName = "script"
	if !e.isAlwaysClose() {
		t.Fatal("isAlwaysClose on script should be true")
	}
}

func TestElement_String(t *testing.T) {
	e := &element{}
	e.nodeName = "script"
	if e.String() != "<script></script>" {
		t.Fatal("Element.String test failed")
	}
}
