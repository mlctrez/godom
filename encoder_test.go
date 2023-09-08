package godom

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncoder_Attributes(t *testing.T) {
	e := &encoder{}
	e.Attributes([]*Attribute{{Name: "foo", Value: "bar"}})
	if len(e.tokens) != 1 {
		t.Fatal("tokens length incorrect")
	}
	attr := e.tokens[0].(*xmlAttr)
	if attr.a[0].Name != "foo" {
		t.Fatal("incorrect attribute name")
	}
	if attr.a[0].Value != "bar" {
		t.Fatal("incorrect attribute wasmValue")
	}
}

func TestEncoder_Indent(t *testing.T) {
	exp := `<html>
  <meta/>
</html>`
	d := Document()
	html := d.CreateElement("html")
	html.AppendChild(d.CreateElement("meta"))

	enc := NewEncoder(&bytes.Buffer{})
	enc.Indent("  ")
	result := html.Marshal(enc).Xml()

	if result != exp {
		t.Fatalf("expected %q but got %q", exp, result)
	}
}

func TestEncoder_AttributesEncoded(t *testing.T) {
	exp := `<html lang="EN"/>`
	d := Document()
	html := d.CreateElement("html")
	html.SetAttribute("lang", "EN")

	enc := NewEncoder(&bytes.Buffer{})
	result := html.Marshal(enc).Xml()

	if result != exp {
		t.Fatalf("expected %q but got %q", exp, result)
	}
}

func TestEncoder_Text(t *testing.T) {
	exp := `<html>some text</html>`
	d := Document()
	html := d.CreateElement("html")
	html.AppendChild(d.CreateTextNode("some text"))

	enc := NewEncoder(&bytes.Buffer{})
	result := html.Marshal(enc).Xml()

	if result != exp {
		t.Fatalf("expected %q but got %q", exp, result)
	}
}

func TestEncoder_lookBack(t *testing.T) {
	e := &encoder{}
	e.tokens = append(e.tokens, &xmlText{d: "1"})
	e.tokens = append(e.tokens, &xmlText{d: "2"})
	e.tokens = append(e.tokens, &xmlText{d: "3"})

	back := e.lookBack(2)
	if back != nil {
		t.Fatal("lookBack test for nil failed")
	}
}

func TestEncoder_shouldClose(t *testing.T) {
	a := assert.New(t)
	api := Document().DocApi()
	a.Equal("<p/>", api.El("p").String())
	a.Equal("<script></script>", api.El("script").String())

}
