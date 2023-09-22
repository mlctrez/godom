package godom

import (
	"bytes"
	"strings"
)

// Element is defined by https://developer.mozilla.org/en-US/docs/Web/API/Element
type Element interface {
	Node
	Remove()
	RemoveChild(child Value)
	SetAttribute(name string, value interface{})
	ReplaceWith(replacement Node)
	SetParent(parent Element)
	Parent() Element
	GetElementsByTagName(tag string) []Element
}

var _ Element = (*element)(nil)

type element struct {
	node
	attributes Attributes
	parent     Element
}

func (e *element) ReplaceWith(n Node) {
	e.this.Call("replaceWith", n.This())
	e.Remove()
}

func (e *element) SetAttribute(name string, value interface{}) {
	e.attributes = append(e.attributes, &Attribute{Name: name, Value: value})
	e.this.Call("setAttribute", name, value)
	e.attributes.SortByName()
}

func (e *element) NodeType() NodeType {
	return NodeTypeElement
}

func (e *element) AppendChild(child Node) {
	e.node.AppendChild(child)
	if el, ok := child.(Element); ok {
		el.SetParent(e)
	}
}

func (e *element) RemoveChild(child Value) {
	var updated []Node
	for _, otherChild := range e.children {
		if otherChild.This().Equal(child) {
			e.this.Call("removeChild", child)
			continue
		}
		updated = append(updated, otherChild)
	}
	e.children = updated
}

func (e *element) Remove() {
	for _, child := range e.children {
		if el, ok := child.(Element); ok {
			el.Remove()
		}
	}
	e.cleanUp()
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
	return e.nodeName == "script"
}

func (e *element) String() string {
	enc := NewEncoder(&bytes.Buffer{})
	return e.Marshal(enc).Xml()
}

func (e *element) SetParent(parent Element) {
	e.parent = parent
}

func (e *element) Parent() Element {
	return e.parent
}

func (e *element) GetElementsByTagName(tag string) []Element {
	var result []Element
	if e.NodeName() == tag {
		result = append(result, e)
	}
	for _, child := range e.children {
		if childEl, ok := child.(*element); ok {
			result = append(result, childEl.GetElementsByTagName(tag)...)
		}
	}
	return result
}

func ElementFromValue(value Value) Element {
	if elType, ok := value.GoValue().(Element); ok {
		return elType
	}

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
