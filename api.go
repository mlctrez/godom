package godom

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// Doc is a helper class for working with the Document interface.
type Doc struct {
	Doc Document
}

// El creates a new element with optional attributes.
func (d Doc) El(tag string, attributes ...*Attribute) Element {
	c := d.Doc.CreateElement(tag)
	for _, a := range attributes {
		c.SetAttribute(a.Name, a.Value)
	}
	return c
}

// At creates a new attribute.
func (d Doc) At(name string, value interface{}) *Attribute {
	return &Attribute{Name: name, Value: value}
}

// T creates a new text with optional attributes.
func (d Doc) T(text string) Text {
	return d.Doc.CreateTextNode(text)
}

func (d Doc) H(html string) Node {
	bufferString := bytes.NewBufferString(html)
	n, err := d.Decode(xml.NewDecoder(bufferString))
	if err != nil {
		errDiv := d.El("div", d.At("style", "color:red;"))
		errDiv.AppendChild(d.T(err.Error()))
		return errDiv
	}
	return n
}

func (d Doc) Decode(decoder *xml.Decoder) (Node, error) {
	var parents []Node
	charBuffer := &charDataBuffer{}

	startNode := func(doc Doc, x xml.StartElement) Element {
		var ga []*Attribute
		for _, attr := range x.Attr {
			ga = append(ga, &Attribute{Name: attr.Name.Local, Value: attr.Value})
		}
		return doc.El(x.Name.Local, ga...)
	}

	token, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("EOF first token")
	}
	for ; err == nil; token, err = decoder.Token() {
		switch x := token.(type) {
		case xml.Directive:
			// for now, we don't care about DOCTYPE, CDATA, etc.
		case xml.StartElement:
			parents = append(parents, startNode(d, x))
			if len(parents) > 1 {
				parents[len(parents)-2].AppendChild(parents[len(parents)-1])
			}
		case xml.EndElement:
			textData := strings.TrimSpace(charBuffer.pop())
			if len(textData) > 0 {
				parents[len(parents)-1].AppendChild(d.T(textData))
			}
			if len(parents) > 1 {
				parents = parents[:len(parents)-1]
			}
		case xml.CharData:
			charBuffer.Write(x)
		}
	}
	if err != io.EOF && err != nil {
		return nil, err
	}
	return parents[0], nil
}

type charDataBuffer struct {
	bytes.Buffer
}

func (cd *charDataBuffer) pop() string {
	charData := cd.String()
	cd.Reset()
	return charData
}

const IM = "implement me"
