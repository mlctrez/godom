//go:build !wasm

package godom

var globalWindow Window

// GlobalClear clears the global for window
func GlobalClear() {
	globalWindow = nil
}

func Global() Window {
	if globalWindow == nil {
		wv := valueT(TypeObject)
		wv.set("document", initialDocument())
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
	e.set("appendChild", func(v Value) {})
	return e
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
	v.set("createTextNode", func() Value {
		return valueT(TypeObject).set("document", d)
	})
	return v
}

var _ Window = (*window)(nil)

type window struct {
	this Value
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
