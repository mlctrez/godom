package godom

import (
	"bytes"
	"fmt"
)

type Encoder interface {
	Start(name string)
	Attributes(a []*Attribute)
	Text(data string)
	End(name string, alwaysClose bool)
	Indent(indentAmount string)
	Flush()
	Xml() string
}

type Marshaller interface {
	Marshal(encoder Encoder) Encoder
}

func NewEncoder(buf *bytes.Buffer) Encoder {
	enc := &encoder{buf: buf}
	return enc
}

var _ Encoder = (*encoder)(nil)

type encoder struct {
	buf          *bytes.Buffer
	tokens       []interface{}
	indentAmount string
}

type xmlStart struct{ n string }
type xmlAttr struct{ a []*Attribute }
type xmlText struct{ d string }
type xmlEnd struct {
	n string
	c bool
}

func (x *encoder) token(t interface{}) { x.tokens = append(x.tokens, t) }

func (x *encoder) Start(name string) { x.token(&xmlStart{n: name}) }

func (x *encoder) Attributes(a []*Attribute) {
	if len(a) > 0 {
		x.token(&xmlAttr{a: a})
	}
}

func (x *encoder) Text(data string) { x.token(&xmlText{d: data}) }

func (x *encoder) End(name string, withCloseElement bool) {
	x.token(&xmlEnd{n: name, c: withCloseElement})
}

func (x *encoder) Indent(indentAmount string) { x.indentAmount = indentAmount }

func (x *encoder) Xml() string {
	x.Flush()
	str := x.buf.String()
	x.buf.Reset()
	return str
}

func (x *encoder) shouldClose(i int) bool {
	if i > 0 {
		switch x.tokens[i-1].(type) {
		case *xmlAttr, *xmlStart:
			return true
		}
	}
	return false
}

func (x *encoder) lookBack(s int) *xmlStart {
	for i := s - 1; i >= 0; i-- {
		if sr, ok := x.tokens[i].(*xmlStart); ok {
			return sr
		}
	}
	return nil
}

func (x *encoder) Flush() {
	indent := -1
	for i, token := range x.tokens {
		switch t := token.(type) {
		case *xmlStart:
			if x.shouldClose(i) {
				x.buf.WriteString(">")
			}
			indent++
			if i > 0 {
				x.buf.WriteString(pad(indent, x.indentAmount))
			}
			x.buf.WriteString("<" + t.n)
		case *xmlAttr:
			for _, i := range t.a {
				x.buf.WriteString(fmt.Sprintf(" %s=%q", i.Name, i.Value))
			}
		case *xmlEnd:
			if x.shouldClose(i) {
				if t.c {
					x.buf.WriteString("></" + t.n + ">")
				} else {
					x.buf.WriteString("/>")
				}
			} else {
				lb := x.lookBack(i)
				if lb != nil && lb.n != t.n {
					x.buf.WriteString(pad(indent, x.indentAmount))
				}
				x.buf.WriteString("</" + t.n + ">")
			}
			indent--
		case *xmlText:
			if x.shouldClose(i) {
				x.buf.WriteString(">")
			}
			x.buf.WriteString(t.d)
		}
	}
}

func pad(indent int, indentAmount string) (result string) {
	if indentAmount == "" {
		return ""
	}
	result = "\n"
	for i := 0; i < indent; i++ {
		result += indentAmount
	}
	return
}
