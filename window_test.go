package godom

import (
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
	w := global.(*window)
	if w.this == nil {
		t.Fatal("window.this was nil")
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

//func TestGlobal_Navigator(t *testing.T) {
//	if Global().Navigator() == nil {
//		t.Fatal("Global().Navigator() returned nil")
//	}
//}
