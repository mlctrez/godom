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
	Api() Doc
}

var _ Document = (*document)(nil)

type document struct {
	node
	documentElement Element
}

func (d *document) Api() Doc {
	return Doc{Doc: d}
}

func (d *document) DocumentElement() Element {
	if d.documentElement != nil {
		fmt.Println("d.documentElement != nil")
		return d.documentElement
	}
	d.documentElement = ElementFromValue(d.this.Get("documentElement"))
	return d.documentElement
}

func (d *document) SetDocumentElement(e Element) {
	d.documentElement = e
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
	return TextFromValue(d.this.Call("createTextNode", data))
}

func (d *document) CreateElement(tagName string) Element {
	return ElementFromValue(d.this.Call("createElement", tagName))
}

func (d *document) CreateElementNS(ns, nodeName string) Element {
	return ElementFromValue(d.this.Call("createElementNS", ns, nodeName))
}

func (d *document) NodeType() NodeType {
	return NodeTypeDocument
}

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
