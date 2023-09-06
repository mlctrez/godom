package godom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestElementFromValue(t *testing.T) {
	a := assert.New(t)
	el := Global().Document().DocumentElement()

	a.Equal(NodeTypeElement, int(el.NodeType()))
	a.Equal("html", el.NodeName())
}
