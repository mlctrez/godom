//go:build !wasm

package godom

var mockGlobal *value

func global() Value {
	if mockGlobal == nil {
		mockGlobal = valueT(TypeObject,
			dataMap{"window": mockWindow()},
		)
	}
	return mockGlobal
}

func mockWindow() Value {
	return valueT(TypeObject,
		dataMap{"document": mockDocument()},
	)
}

func mockDocument() *value {
	doc := valueT(TypeObject)
	doc.data = dataMap{
		"createElement": func(nodeName string) Value {
			return elementValue(doc, "", nodeName)
		},
		"createElementNS": func(ns, nodeName string) Value {
			return elementValue(doc, ns, nodeName)
		},
		"createTextNode": func(text string) Value {
			return valueT(TypeObject).set("document", doc).set("data", text)
		},
	}
	return doc
}

func elementValue(doc Value, ns, nodeName string) *value {
	e := valueT(TypeObject)
	e.data = dataMap{
		"namespaceURI":  ns,
		"nodeName":      nodeName,
		"document":      doc,
		"remove":        func() { panic(IM) },
		"replaceWith":   func(v Value) { panic(IM) },
		"hasAttributes": func() Value { return ToValue(false) },
		"hasChildNodes": func() Value { return ToValue(false) },
		"attributes":    func() Value { return ToValue([]interface{}{}) },
		"childNodes":    func() Value { return ToValue([]interface{}{}) },
		"appendChild":   func(v Value) {},
		"setAttribute":  func(name string, value interface{}) {},
	}
	return e
}
