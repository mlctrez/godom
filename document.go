package godom

import "fmt"

// Document is defined by https://developer.mozilla.org/en-US/docs/Web/API/Document
type Document interface {
	Node
	DocumentElement() Element
	SetDocumentElement(Element)
	Head() Element
	Body() Element
	CreateTextNode(text string) Text
	CreateElement(tagName string) Element
	CreateElementNS(namespaceURI, qualifiedName string) Element
}

var _ Document = (*document)(nil)

type document struct {
	node
}

func (d *document) DocumentElement() Element {
	return d.this.Get("documentElement").GoValue().(Element)
}

func (d *document) SetDocumentElement(e Element) {
	d.this.Set("documentElement", e.This())
}

func (d *document) Head() Element {
	nodes := d.DocumentElement().ChildNodes()
	return findElement(nodes, "head", true)
}

func (d *document) Body() Element {
	nodes := d.DocumentElement().ChildNodes()
	return findElement(nodes, "body", true)
}

func (d *document) CreateTextNode(data string) Text {
	t := &text{data: data}
	t.data = data
	t.this = d.this.Call("createTextNode")
	t.this.SetGoValue(t)
	return t
}

func (d *document) CreateElement(tagName string) Element {
	e := &element{}
	e.nodeName = tagName
	e.this = d.this.Call("createElement", tagName)
	e.this.SetGoValue(e)
	return e
}

func (d *document) CreateElementNS(ns, nodeName string) Element {
	e := &element{}
	e.ns = ns
	e.nodeName = nodeName
	e.this = d.this.Call("createElementNS", ns, nodeName)
	e.this.SetGoValue(e)
	return e
}

func (d *document) NodeType() NodeType { return NodeTypeDocument }

func findElement(children []Node, nodeName string, panicOnError bool) Element {
	for _, child := range children {
		if child.NodeName() == nodeName && child.NodeType() == NodeTypeElement {
			return child.(Element)
		}
	}
	if panicOnError {
		panic(fmt.Errorf("unable to find Element %q", nodeName))
	}
	return nil
}
