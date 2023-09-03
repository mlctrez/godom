package godom

import "strings"

// Text is described by https://developer.mozilla.org/en-US/docs/Web/API/Text
type Text interface {
	Node
	String() string
	IsWhiteSpace() bool
}

var _ Text = (*text)(nil)

type text struct {
	node
	data string
}

func (t *text) Marshal(enc Encoder) Encoder {
	enc.Text(t.data)
	return enc
}

func (t *text) String() string {
	return t.data
}

func (t *text) IsWhiteSpace() bool {
	return strings.TrimSpace(t.data) == ""
}

func (t *text) NodeType() NodeType {
	return NodeTypeText
}

func TextFromValue(value Value) Text {
	t := &text{}
	t.this = value
	t.data = value.Get("data").String()
	value.SetGoValue(t)
	return t
}
