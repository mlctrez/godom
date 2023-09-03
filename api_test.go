package godom

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
)

func TestDoc_El(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
	html := doc.El("html", doc.At("lang", "en"))
	if html.NodeName() != "html" {
		t.Fatal("html node name not set correctly")
	}
}

func TestDoc_At(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
	at := doc.At("name", "wasmValue")
	if at.Name != "name" {
		t.Fatal("name not set correctly")
	}
	if at.Value != "wasmValue" {
		t.Fatal("wasmValue not set correctly")
	}
}

func TestDoc_FromDecoder(t *testing.T) {
	doc := Doc{Doc: Global().Document()}
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

func testEncodeHelper(n Node) string {
	buf := &bytes.Buffer{}
	enc := NewEncoder(buf)
	n.Marshal(enc).Flush()

	return buf.String()
}
