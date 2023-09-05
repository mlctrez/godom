package godom

import (
	"runtime"
	"testing"
)

func TestGlobal(t *testing.T) {
	global := Global()
	if global == nil {
		t.Fatal("Global() returned nil")
	}
	if ptr(global) != ptr(Global()) {
		t.Fatal("multiple calls returned different object")
	}
}

func TestGlobal_Document(t *testing.T) {
	g := Global()
	d := g.Document()
	if d == nil {
		t.Fatal("Global().Document() returned nil")
	}
	if ptr(d) != ptr(g.Document()) {
		t.Fatal("multiple calls returned different object")
	}

}

func TestLocation_Href(t *testing.T) {
	href := Global().Location().Href()
	if runtime.GOOS == "linux" {
		if "http://testserver" != href {
			t.Error("unexpected href")
		}
	}
}
