package godom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNode_AddEventListener(t *testing.T) {
	a := assert.New(t)
	a.NotNil(a)

	var eventHappened Value

	d := Global().Document()
	release := d.AddEventListener("click", func(event Value) {
		eventHappened = event
	})
	defer release()

	d.This().Call("dispatchEvent", Global().Value().Get("Event").New("click"))
	a.NotNil(eventHappened)

	a.Equal("click", eventHappened.Get("type").String())
	a.True(d.This().Equal(eventHappened.Get("target")))
}
