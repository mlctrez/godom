//go:build !wasm

package godom

//func Test_mockDocument(t *testing.T) {
//	a := assert.New(t)
//	v := mockDocument()
//	result := v.Call("createElementNS", "a", "b")
//	a.Equal("b", result.Get("nodeName").String())
//}
//
//func Test_elementValue(t *testing.T) {
//	a := assert.New(t)
//	v := elementValue(mockDocument(), "", "div")
//
//	testPanic := func(test string, target func()) {
//		defer func() { a.NotNilf(recover(), "test %s failed", test) }()
//		target()
//	}
//	testPanic("remove", func() { v.Call("remove") })
//	testPanic("replaceWith", func() { v.Call("replaceWith", ToValue("")) })
//
//	a.True(v.Call("attributes").Length() == 0)
//	a.True(v.Call("childNodes").Length() == 0)
//}
