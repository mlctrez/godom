package godom

import (
	"reflect"
	"testing"
)

func TestDocument_Body(t *testing.T) {
	globalClear()
	doc := Doc{Doc: Global().Document()}
	doc.Doc.SetDocumentElement(doc.El("html"))
	root := doc.Doc.DocumentElement()
	root.AppendChild(doc.El("body"))
	if "body" != doc.Doc.Body().NodeName() {
		t.Fatal("node not found")
	}
}

func TestDocument_Head(t *testing.T) {
	globalClear()
	doc := Doc{Doc: Global().Document()}
	doc.Doc.SetDocumentElement(doc.El("html"))
	root := doc.Doc.DocumentElement()
	root.AppendChild(doc.El("head"))
	if "head" != doc.Doc.Head().NodeName() {
		t.Fatal("node not found")
	}
}

func TestDocument_CreateTextNode(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
	textNode := doc.Doc.CreateTextNode("some text")
	if reflect.TypeOf(textNode) != reflect.TypeOf(&text{}) {
		t.Fatal("incorrect node type")
	}
	if textNode.(*text).data != "some text" {
		t.Fatal("incorrect data")
	}
}

func TestDocument_CreateElementNS(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
	ns := "http://www.w3.org/2000/svg"
	svg := doc.Doc.CreateElementNS(ns, "svg")
	if svg.(*element).ns != ns {
		t.Fatal("namespace not correct")
	}
}

func TestDocument_NodeType(t *testing.T) {
	nt := (Global().Document()).NodeType()
	if nt != NodeTypeDocument {
		t.Fatal("incorrect node type")
	}
}

func TestDocument_DocumentElement(t *testing.T) {
	Global().Document().DocumentElement()
}

func TestDocument_findElement_fail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("findElement did not panic")
		}
	}()
	children := []Node{&node{nodeName: "foo"}}
	findElement(children, "bar", true)
}

func TestDocument_findElement_nopanic(t *testing.T) {
	children := []Node{&node{nodeName: "foo"}}
	e := findElement(children, "bar", false)
	if e != nil {
		t.Fatal("nil element not returned")
	}
}
