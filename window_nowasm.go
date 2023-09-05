//go:build !wasm

package godom

import (
	"fmt"
	"log"
)

var globalWindow Window

func globalClear() {
	globalWindow = nil
}

func Global() Window {
	if globalWindow == nil {
		wv := valueT(TypeObject)
		wv.set("document", initialDocument())
		wv.set("console", initialConsole())
		wv.set("location", initialLocation())
		globalWindow = &window{this: wv}
		wv.SetGoValue(globalWindow)
	}
	return globalWindow
}

func elementValue(d *document, ns, nodeName string) *value {
	e := valueT(TypeObject)
	e.set("namespaceURI", ns)
	e.set("nodeName", nodeName)
	e.set("document", d)
	e.set("remove", func() {})
	e.set("replaceWith", func(v Value) {})
	e.set("hasAttributes", func() Value { return ToValue(false) })
	e.set("hasChildNodes", func() Value { return ToValue(false) })
	e.set("attributes", func() Value { return ToValue([]interface{}{}) })
	e.set("childNodes", func() Value { return ToValue([]interface{}{}) })
	e.set("appendChild", func(v Value) {})
	e.set("setAttribute", func(name string, value interface{}) {})
	return e
}

func initialConsole() Value {
	c := &console{}
	v := valueT(TypeObject)
	c.this = v
	v.SetGoValue(c)
	v.set("log", func(args ...interface{}) {
		log.Println(args...)
	})
	v.set("error", func(args ...interface{}) {
		log.Println(args...)
	})
	return v
}

func initialDocument() Value {
	d := &document{}
	v := valueT(TypeObject)
	d.this = v
	v.SetGoValue(d)
	v.set("createElement", func(tagName string) Value {
		return elementValue(d, "", tagName)
	})
	v.set("createElementNS", func(ns, nodeName string) Value {
		return elementValue(d, ns, nodeName)
	})
	v.set("createTextNode", func(data string) Value {
		return valueT(TypeObject).set("document", d).set("data", data)
	})
	return v
}

func initialLocation() Value {
	d := &location{}
	v := valueT(TypeObject)
	d.this = v
	v.SetGoValue(d)
	v.set("href", "http://testserver")
	v.set("reload", func() {
		fmt.Println("reload called")
	})
	return v
}

var _ Window = (*window)(nil)

type window struct {
	this Value
}

func (w *window) Value() Value {
	return w.this
}

func (w *window) value() *value {
	return w.this.(*value)
}

func (w *window) Document() Document {
	return w.value().data["document"].(Value).GoValue().(Document)
}

func (w *window) Navigator() Navigator {
	return w.value().data["navigator"].(Value).GoValue().(Navigator)
}

func (w *window) Location() Location {
	return w.value().data["location"].(Value).GoValue().(Location)
}

func (w *window) Console() Console {
	return w.value().data["console"].(Value).GoValue().(Console)
}
