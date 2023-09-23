package godom

import (
	"fmt"
)

func Global() Value {
	return global()
}

type DocumentApi interface {
	Node
	CreateElement(tag string) Element
	CreateElementNS(tag string, ns string) Element
	CreateTextNode(text string) Text
	DocumentElement() Element
	Head() Element
	Body() Element
	EventListener
	DocApi() DocApi
}

type document struct {
	node
}

func (d *document) DocApi() DocApi {
	return NewDocApi(d.this)
}

func (d *document) CreateElement(tag string) Element {
	return ElementFromValue(d.this.Call("createElement", tag))
}

func (d *document) CreateElementNS(tag string, ns string) Element {
	return ElementFromValue(d.this.Call("createElementNS", tag, ns))
}

func (d *document) CreateTextNode(text string) Text {
	return TextFromValue(d.this.Call("createTextNode", text))
}

func (d *document) DocumentElement() Element {
	// TODO: appropriate caching
	return ElementFromValue(d.this.Get("documentElement"))
}

func (d *document) Head() (body Element) {
	return d.child("head")
}

func (d *document) Body() (body Element) {
	return d.child("body")
}

func (d *document) child(tag string) (body Element) {
	for _, child := range d.DocumentElement().ChildNodes() {
		if child.NodeName() == tag {
			return child.(*element)
		}
	}
	panic(fmt.Sprintf("unable to locate %s", tag))
}

func Document() DocumentApi {
	docElement := Global().Get("window").Get("document")
	d := &document{node{this: docElement}}
	d.node.children = d.DocumentElement().ChildNodes()
	return d
}
