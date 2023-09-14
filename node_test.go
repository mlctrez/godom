package godom

import (
	"bytes"
	"testing"
)

func TestNode_Marshal(t *testing.T) {
	enc := NewEncoder(&bytes.Buffer{})
	n := &node{nodeName: "node"}
	xml := n.Marshal(enc).Xml()
	if "<node/>" != xml {
		t.Fatal("expected <node/> but got" + xml)
	}
}

func TestNode_NodeType(t *testing.T) {
	n := &node{nodeName: "node"}
	if NodeTypeNone != n.NodeType() {
		t.Fatalf("expected NodeTypeNone but got %d", n.NodeType())
	}
}

func TestNode_AppendChild(t *testing.T) {

	doc := Document().DocApi()
	e := doc.El("html")
	e.AppendChild(doc.El("body"))
	if len(e.ChildNodes()) != 1 {
		t.Fatal("AppendChild did not result in one child")
	}
	if e.ChildNodes()[0].NodeName() != "body" {
		t.Fatal("AppendChild did not result in correct child")
	}
}

func TestNode_String(t *testing.T) {
	n := &node{nodeName: "node"}
	if n.String() != "Node:node" {
		t.Fatalf("expected Node:node but got %q", n.String())
	}
}
