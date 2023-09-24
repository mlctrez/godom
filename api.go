package godom

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type OnElement func(e Element, name, data string)

type DocApi interface {
	El(tag string, attributes ...*Attribute) Element
	At(name string, value interface{}) *Attribute
	T(text string) Text
	H(html string) Element

	WithCallback(oe OnElement) DocApi
	Decode(decoder *xml.Decoder) (Element, error)
}

func NewDocApi(v Value) DocApi {
	return &doc{v: v}
}

type doc struct {
	v  Value
	cb OnElement
}

var dataGoRegex = regexp.MustCompile("^data-go.*?|go.*?")

func (d *doc) WithCallback(cb OnElement) DocApi {
	return &doc{v: d.v, cb: cb}
}

// El creates a new element with optional attributes.
func (d *doc) El(tag string, attributes ...*Attribute) Element {
	c := ElementFromValue(d.v.Call("createElement", tag))
	var nv [][]string
	for _, a := range attributes {
		c.SetAttribute(a.Name, a.Value)
		if d.cb != nil && dataGoRegex.MatchString(a.Name) {
			name := strings.TrimPrefix(a.Name, "data-")
			// use later after all attributes have been set
			nv = append(nv, []string{name, a.Value.(string)})
		}
	}
	if d.cb != nil {
		for _, nameValue := range nv {
			d.cb(c, nameValue[0], nameValue[1])
		}
	}
	return c
}

// At creates a new attribute.
func (d *doc) At(name string, value interface{}) *Attribute {
	return &Attribute{Name: name, Value: value}
}

// T creates a new text with optional attributes.
func (d *doc) T(text string) Text {
	return TextFromValue(d.v.Call("createTextNode", text))
}

func (d *doc) H(html string) Element {
	bufferString := bytes.NewBufferString(html)
	n, err := d.Decode(xml.NewDecoder(bufferString))
	if err != nil {
		errDiv := d.El("div", d.At("style", "color:red;"))
		errDiv.AppendChild(d.T(err.Error()))
		return errDiv
	}
	return n
}

func (d *doc) Decode(decoder *xml.Decoder) (Element, error) {
	var parents []Element
	charBuffer := &charDataBuffer{}

	token, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("EOF first token")
	}
	for ; err == nil; token, err = decoder.Token() {
		switch x := token.(type) {
		case xml.Directive:
			// for now, we don't care about DOCTYPE, CDATA, etc.
		case xml.StartElement:
			parents = append(parents, d.startNode(x))
			if len(parents) > 1 {
				d.AppendChild(parents[len(parents)-2], parents[len(parents)-1])
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
	if len(parents) == 0 {
		err = errors.New("no parent elements")
	}
	if err != io.EOF && err != nil {
		return nil, err
	}
	return parents[0], nil
}

func (d *doc) startNode(x xml.StartElement) Element {
	var ga []*Attribute
	for _, attr := range x.Attr {
		ga = append(ga, &Attribute{Name: attr.Name.Local, Value: attr.Value})
	}
	return d.El(x.Name.Local, ga...)
}

func (d *doc) AppendChild(parent, child Element) {
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

/*

 godom.Ser(callback).H("some html here") - returns element

 d := godom.Ser(callback)
 d.El("body").Body(nodes ...Node)



*/
