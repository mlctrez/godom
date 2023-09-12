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
	Doc      Value
	CallBack func(e Element, dataGo []string)
}

// El creates a new element with optional attributes.
func (d Doc) El(tag string, attributes ...*Attribute) Element {
	c := ElementFromValue(d.Doc.Call("createElement", tag))
	var dataGo []string
	for _, a := range attributes {
		if a.Name == "data-go" {
			dataGo = append(dataGo, a.Value.(string))
		} else {
			c.SetAttribute(a.Name, a.Value)
		}
	}
	if dataGo != nil {
		d.CallBack(c, dataGo)
	}
	return c
}

// At creates a new attribute.
func (d Doc) At(name string, value interface{}) *Attribute {
	return &Attribute{Name: name, Value: value}
}

// T creates a new text with optional attributes.
func (d Doc) T(text string) Text {
	return TextFromValue(d.Doc.Call("createTextNode", text))
}

func (d Doc) H(html string) Element {
	bufferString := bytes.NewBufferString(html)
	n, err := d.Decode(xml.NewDecoder(bufferString))
	if err != nil {
		errDiv := d.El("div", d.At("style", "color:red;"))
		errDiv.AppendChild(d.T(err.Error()))
		return errDiv
	}
	return n
}

func (d Doc) Decode(decoder *xml.Decoder) (Element, error) {
	var parents []Element
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
				appendChild(parents[len(parents)-2], parents[len(parents)-1])
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

func appendChild(parent, child Element) {
	parent.AppendChild(child)
	child.SetParent(parent)
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
