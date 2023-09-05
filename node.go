package godom

import (
	"fmt"
)

// Node is defined by https://developer.mozilla.org/en-US/docs/Web/API/Node
type Node interface {
	Marshaller
	EventListener
	// ChildNodes implements Node.childNodes
	ChildNodes() []Node
	NodeName() string
	NodeType() NodeType
	AppendChild(child Node)
	Marshal(enc Encoder) Encoder
	String() string
	This() Value
}

var _ Node = (*node)(nil)

type node struct {
	this     Value
	nodeName string
	ns       string
	children []Node
}

func (d *node) This() Value        { return d.this }
func (d *node) ChildNodes() []Node { return d.children }
func (d *node) NodeName() string   { return d.nodeName }
func (d *node) NodeType() NodeType { return NodeTypeNone }

func (d *node) AppendChild(child Node) {
	d.this.Call("appendChild", child.This())
	d.children = append(d.children, child)
}

func (d *node) Marshal(enc Encoder) Encoder {
	enc.Start(d.nodeName)
	enc.End(d.nodeName, false)
	return enc
}

func (d *node) String() string {
	return fmt.Sprintf("Node:%s", d.nodeName)
}

type NodeType uint

const (
	NodeTypeNone                  = 0
	NodeTypeElement               = 1
	NodeTypeAttribute             = 2
	NodeTypeText                  = 3
	NodeTypeCDATA                 = 4
	NodeTypeProcessingInstruction = 7
	NodeTypeComment               = 8
	NodeTypeDocument              = 9
	NodeTypeDocumentType          = 10
	NodeTypeDocumentFragment      = 11
)
