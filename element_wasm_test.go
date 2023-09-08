package godom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElementFromValue(t *testing.T) {
	a := assert.New(t)
	el := Document().DocumentElement()

	a.Equal(NodeTypeElement, int(el.NodeType()))
	a.Equal("html", el.NodeName())
}

func TestElement_Remove(t *testing.T) {

	name := t.Name()

	a := assert.New(t)
	d := Document()
	div := d.DocApi().El("div")
	div.SetAttribute("id", name)

	body := d.Body()
	body.AppendChild(div)

	testingDiv := d.This().Call("getElementById", name)
	a.True(!testingDiv.IsUndefined())

	div.Remove()

	testingDiv = d.This().Call("getElementById", name)
	a.True(testingDiv.IsUndefined())

}

func TestElement_ReplaceWith(t *testing.T) {
	a := assert.New(t)
	d := Document()
	divOne := d.DocApi().El("div")
	divOne.SetAttribute("id", "testingDivOne")

	divTwo := d.DocApi().El("div")
	divOne.SetAttribute("id", "testingDivTwo")

	body := d.Body()
	body.AppendChild(divOne)

	a.True(!d.This().Call("getElementById", "testingDivOne").IsUndefined())

	divOne.ReplaceWith(divTwo)
	a.True(d.This().Call("getElementById", "testingDivOne").IsUndefined())
	a.True(!d.This().Call("getElementById", "testingDivTwo").IsUndefined())

}
