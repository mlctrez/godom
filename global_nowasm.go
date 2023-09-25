//go:build !wasm

package godom

import (
	"sync"
)

var globalThisVar Value
var globalThisMu sync.Mutex

type mockObject struct {
	props      map[string]interface{}
	attributes map[string]interface{}
}

func (mo *mockObject) init() {
	if mo.props == nil {
		mo.props = map[string]interface{}{}
	}
	if mo.attributes == nil {
		mo.attributes = map[string]interface{}{}
	}
}

func (mo *mockObject) Set(prop string, value any) {
	mo.init()
	mo.props[prop] = value
}

func (mo *mockObject) Get(prop string) Value {
	mo.init()
	return toValue(mo.props[prop])
}

func (mo *mockObject) SetAttribute(prop string, value any) {
	mo.init()
	mo.attributes[prop] = value
}

func (mo *mockObject) GetAttribute(prop string) Value {
	mo.init()
	return toValue(mo.attributes[prop])
}

type globalThis struct {
	mockObject
}

func (t *globalThis) Call(method string, args []interface{}) (Value, bool) {
	if v, ok := t.props[method]; ok {
		if nwf, okF := v.(*noWasmFunc); okF {
			return toValue(nwf.fn(nwf.Value, toValues(args))), true
		}
	}
	return nil, false
}

type mockWindow struct {
	mockObject
}

type mockDocument struct {
	mockObject
}

func (md *mockDocument) CreateElement(tag string) Value {
	me := &mockElement{}
	me.Set("nodeName", tag)
	return toValue(me)
}

func (md *mockDocument) CreateTextNode(text string) Value {
	me := &mockElement{}
	me.Set("data", text)
	return toValue(me)
}

func (md *mockDocument) Reset() {
	documentElement := NewDocApi(toValue(md)).H("<html><head/><body/></html>")
	md.Set("documentElement", documentElement.This())
}

func (md *mockDocument) GetElementById(id string) Value {
	for _, child := range Document().Body().ChildNodes() {
		if child.This().Call("getAttribute", "id").String() == id {
			return child.This()
		}
	}
	return toValue(nil)
}

type mockElement struct {
	mockObject
	children []Value
}

func (m *mockElement) AppendChild(child Value) {
	m.children = append(m.children, child)
}

func (m *mockElement) RemoveChild(child Value) {
	var updated []Value
	for _, otherChild := range m.children {
		if otherChild.Equal(child) {
			continue
		}
		updated = append(updated, otherChild)
	}
	m.children = updated
}

var dummyFunc = func() {}

func (m *mockElement) AddEventListener(eventType string, fn any) func() { return dummyFunc }

func (m *mockElement) RemoveEventListener(eventType string, fn any) {}

func (m *mockElement) Remove() {}

func global() Value {
	globalThisMu.Lock()
	defer globalThisMu.Unlock()
	if globalThisVar == nil {
		globalThisVar = toValue(&globalThis{})
		window := toValue(&mockWindow{})

		mockDoc := &mockDocument{}
		mockDoc.Reset()
		window.Set("document", toValue(mockDoc))
		globalThisVar.Set("window", window)

		globalThisVar.Set("console", toValue(&consoleValue{}))

		globalThisVar.Set("Object", &value{t: TypeFunction, v: &ValueObject{}})
		globalThisVar.Set("NaN", &value{t: TypeNumber, v: ValueNaN{}})
	}
	return globalThisVar
}
