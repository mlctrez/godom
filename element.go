package godom

import (
	"bytes"
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
	panic("implement me")
}

func (e *element) SetAttribute(name string, value interface{}) {
	e.attributes = append(e.attributes, &Attribute{Name: name, Value: value})
	e.attributes.SortByName()
}

func (e *element) NodeType() NodeType {
	return NodeTypeElement
}

func (e *element) Remove() {
	panic(IM)
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
