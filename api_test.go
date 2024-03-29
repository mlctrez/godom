package godom

import (
	"bytes"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestDoc_El(t *testing.T) {
	a := assert.New(t)
	doc := Document().DocApi()
	html := doc.El("html", doc.At("lang", "en"))
	a.Equal("html", html.NodeName())
}

func TestDoc_At(t *testing.T) {
	doc := Document().DocApi()
	at := doc.At("name", "wasmValue")
	if at.Name != "name" {
		t.Fatal("name not set correctly")
	}
	if at.Value != "wasmValue" {
		t.Fatal("wasmValue not set correctly")
	}
}

func TestDoc_FromDecoder(t *testing.T) {
	doc := Document().DocApi()
	decodeString := func(s string) (Node, error) {
		return doc.Decode(xml.NewDecoder(bytes.NewBufferString(s)))
	}
	var n Node
	var err error
	n, err = decodeString("")
	if err.Error() != "EOF first token" {
		t.Fatal("bad initial eof handling")
	}

	n, err = decodeString("<asdf>")
	if !strings.Contains(err.Error(), "unexpected EOF") {
		t.Fatal("xml syntax error not passed back")
	}

	n, err = decodeString("<html>data</html>")
	if err != nil {
		t.Fatal()
	}
	if n.NodeName() != "html" {
		t.Fatal()
	}
	enc := testEncodeHelper(n)
	if enc != "<html>data</html>" {
		t.Fatal(enc)
	}
}

func TestDoc_H(t *testing.T) {
	req := require.New(t)
	doc := Document().DocApi()
	html := `<div><button id="one">button text</button></div>`

	req.Equal(html, doc.H(html).String())
	req.Equal("<div style=\"color:red;\">EOF first token</div>", doc.H("").String())
}

func testEncodeHelper(n Node) string {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)
	n.Marshal(enc).Flush()

	return buf.String()
}

func TestDoc_directive(t *testing.T) {

	doc := Document().DocApi()
	decodeString := func(s string) (Node, error) {
		return doc.Decode(xml.NewDecoder(bytes.NewBufferString(s)))
	}
	_, err := decodeString("<!DOCTYPE html><html/>")
	if err != nil {
		t.Fatal(err)
	}

}

func TestDoc_noParent(t *testing.T) {
	a := require.New(t)
	doc := Document().DocApi()
	decodeString := func(s string) (Node, error) {
		return doc.Decode(xml.NewDecoder(bytes.NewBufferString(s)))
	}
	_, err := decodeString("foo")
	a.NotNil(err)
	a.ErrorContains(err, "no parent element")

}

func TestDoc_El_datago(t *testing.T) {
	a := assert.New(t)
	api := Document().DocApi()
	var dataGoRef string
	api = api.WithCallback(func(e Element, name, value string) {
		a.Equal("div", e.NodeName())
		e.SetAttribute("other", "new")
		dataGoRef = value
	})
	h := api.H(`<div data-go="data-go-value" other="old"/>`)
	a.Equal("data-go-value", dataGoRef)

	// callback should have replaced the other attribute with new
	// don't currently have getAttribute so just use string
	a.Equal(`<div data-go="data-go-value" other="new"/>`, h.String())
}

func Test_dataGoRegex(t *testing.T) {
	a := require.New(t)
	r := dataGoRegex
	a.False(r.MatchString(""))
	a.False(r.MatchString("foo"))
	a.True(r.MatchString("data-go"))
	a.True(r.MatchString("data-go-extra"))
	a.True(r.MatchString("go"))
	a.True(r.MatchString("go-extra"))
}

func Test_nested_with_text(t *testing.T) {
	a := require.New(t)
	doc := Document().DocApi()
	a.Equal(
		"<p>some text <code>markup</code> other text</p>",
		doc.H("<p>some text <code>markup</code> other text</p>").String(),
	)
}
