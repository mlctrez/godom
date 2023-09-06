package godom

import (
	"bytes"
	"strings"
)

// Element is defined by https://developer.mozilla.org/en-US/docs/Web/API/Element
type Element interface {
	Node
	Remove()
	SetAttribute(name string, value interface{})
	ReplaceWith(replacement Node)
}

var _ Element = (*element)(nil)

type element struct {
	node
	attributes Attributes
}

func (e *element) ReplaceWith(n Node) {
	e.this.Call("replaceWith", n.This())
}

func (e *element) SetAttribute(name string, value interface{}) {
	if strings.HasPrefix(name, "on") {
		e.this.Set(name, value)
		return
	}
	e.attributes = append(e.attributes, &Attribute{Name: name, Value: value})
	e.this.Call("setAttribute", name, value)
	e.attributes.SortByName()
}

func (e *element) NodeType() NodeType {
	return NodeTypeElement
}

func (e *element) Remove() {
	e.this.Call("remove")
}

func (e *element) Marshal(enc Encoder) Encoder {
	enc.Start(e.nodeName)
	enc.Attributes(e.attributes)
	for _, child := range e.children {
		child.Marshal(enc)
	}
	enc.End(e.nodeName, e.isAlwaysClose())
	return enc
}

func (e *element) isAlwaysClose() bool {
	if e.nodeName == "script" {
		return true
	}
	return false
}

func (e *element) String() string {
	enc := NewEncoder(&bytes.Buffer{})
	return e.Marshal(enc).Xml()
}

func ElementFromValue(value Value) Element {

	e := &element{}
	e.this = value
	nodeName := value.Get("nodeName").String()
	e.nodeName = strings.ToLower(nodeName)
	if !value.Get("namespaceURI").IsNull() {
		e.ns = value.Get("namespaceURI").String()
	}
	value.SetGoValue(e)

	if value.Get("hasAttributes").Truthy() && value.Call("hasAttributes").Bool() {
		attributes := value.Get("attributes")
		for i := 0; i < attributes.Length(); i++ {
			attribute := attributes.Index(i)
			e.attributes = append(e.attributes,
				&Attribute{Name: attribute.Get("name").String(), Value: attribute.Get("value")},
			)
		}
	}

	if value.Get("hasChildNodes").Truthy() && value.Call("hasChildNodes").Bool() {
		children := value.Get("childNodes")
		for i := 0; i < children.Length(); i++ {
			child := children.Index(i)
			switch child.Get("nodeType").Int() {
			case NodeTypeElement:
				elementChild := ElementFromValue(child)
				e.children = append(e.children, elementChild)
			case NodeTypeText:
				textChild := TextFromValue(child)
				if !textChild.IsWhiteSpace() {
					e.children = append(e.children, textChild)
				}
			}
		}
	}

	return e
}
