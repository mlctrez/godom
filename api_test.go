package godom

import (
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
