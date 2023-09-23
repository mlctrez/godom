package godom

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestElementFromValue(t *testing.T) {
	a := assert.New(t)
	el := Document().DocumentElement()

	a.Equal(NodeTypeElement, int(el.NodeType()))
	a.Equal("html", el.NodeName())
}

func TestDocument_CreateElementNS(t *testing.T) {
	a := assert.New(t)
	el := Document().CreateElementNS("someNamespace", "div")
	el.SetAttribute("id", t.Name())
	a.Equal(NodeTypeElement, int(el.NodeType()))
	a.Equal("div", el.NodeName())
}

func TestElement_Remove(t *testing.T) {

	name := t.Name()

	a := assert.New(t)
	d := Document()
	div := d.DocApi().El("div", &Attribute{Name: "id", Value: name})
	div.AddEventListener("eventName", func(event Value) {})
	div.AppendChild(d.DocApi().El("p"))

	isNull := func(elementId string) bool {
		return d.This().Call("getElementById", elementId).IsNull()
	}

	a.True(isNull(name))
	body := d.Body()
	body.AppendChild(div)
	a.False(isNull(name))

	body.RemoveChild(div.This())
	a.True(isNull(name))

	body.AppendChild(div)
	a.False(isNull(name))
	div.Remove()
	a.True(isNull(name))

}

func TestElement_ReplaceWith(t *testing.T) {
	a := assert.New(t)
	d := Document()
	divOne := d.DocApi().El("div")
	divOneId := t.Name() + "One"
	divOne.SetAttribute("id", divOneId)

	divTwo := d.DocApi().El("div")
	divTwoId := t.Name() + "Two"
	divTwo.SetAttribute("id", divTwoId)

	isNull := func(elementId string) bool {
		return d.This().Call("getElementById", elementId).IsNull()
	}

	a.True(isNull(divOneId))
	a.True(isNull(divTwoId))
	body := d.Body()
	body.AppendChild(divOne)

	a.True(!isNull(divOneId))
	a.True(isNull(divTwoId))

	divOne.ReplaceWith(divTwo)

	a.True(isNull(divOneId))
	a.True(!isNull(divTwoId))

}

func TestElement_Parent(t *testing.T) {
	a := assert.New(t)
	doc := Document().DocApi()
	elem := doc.H("<div><p/></div>")
	pElem := elem.ChildNodes()[0].(*element)
	a.Equal("div", pElem.Parent().NodeName())

	elem = doc.H("<div/>")
	elem.AppendChild(doc.H("<p/>"))
	pElem = elem.ChildNodes()[0].(*element)
	a.Equal("div", pElem.Parent().NodeName())
}

func TestElement_GetElementsByTagName(t *testing.T) {
	a := require.New(t)
	h := Document().DocApi().H("<div><span/><span/><div><span/></div></div>")
	elems := h.GetElementsByTagName("span")
	a.Equal(3, len(elems))

	h = Document().DocApi().H("<div><p/><p/><div><p/></div></div>")
	elems = h.GetElementsByTagName("p")
	a.Equal(3, len(elems))
}

func TestElement_RemoveAttribute(t *testing.T) {
	a := require.New(t)
	doc := Document().DocApi()
	h := doc.H(`<div id="removeThisId"/>`)
	a.False(h.This().Call("getAttribute", "id").IsNull())

	a.Equal("removeThisId", h.This().Call("getAttribute", "id").String())
	h.RemoveAttribute("id")
	a.True(h.This().Call("getAttribute", "id").IsNull())
}

func TestElement_Body(t *testing.T) {
	a := require.New(t)
	doc := Document().DocApi()
	h := doc.H(`<div id="testBody"/>`)
	h.Body(doc.H(`<div id="bodyNode1"/>`), doc.H(`<div id="bodyNode2"/>`))
	a.Equal(2, len(h.ChildNodes()))
	a.Equal("bodyNode1", h.ChildNodes()[0].This().Call("getAttribute", "id").String())
}
