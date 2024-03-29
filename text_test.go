package godom

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestText_Marshal(t *testing.T) {
	n := &text{data: "data"}
	enc := NewEncoder(&bytes.Buffer{})
	n.Marshal(enc)
	tokens := enc.(*encoder).tokens
	if len(tokens) != 1 {
		t.Fatal("unexpected tokens length")
	}
	token := tokens[0]
	xt := token.(*xmlText)
	if xt.d != "data" {
		t.Fatal("unexpected token")
	}
}
func TestText_NodeType(t *testing.T) {
	n := &text{data: "data"}
	if n.NodeType() != NodeTypeText {
		t.Fatal("NodeType() test failed")
	}

}
func TestText_String(t *testing.T) {
	s := (&text{data: "data"}).String()
	if s != "data" {
		t.Fatal("String() test failed")
	}
}
func TestText_IsWhiteSpace(t *testing.T) {
	a := assert.New(t)
	a.False((&text{data: "data"}).IsWhiteSpace())
	a.True((&text{data: ""}).IsWhiteSpace())
	a.True((&text{data: "  \n  "}).IsWhiteSpace())
}
