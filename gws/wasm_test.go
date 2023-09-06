//go:build js && wasm

package gws

import (
	"errors"
	"github.com/mlctrez/godom"
	"github.com/stretchr/testify/assert"
	"syscall/js"
	"testing"
)

func Test_handleJsError(t *testing.T) {
	a := assert.New(t)
	defer func() {
		r := recover()
		switch e := r.(type) {
		case error:
			a.Equal("other error", e.Error())
			return
		}
		a.Fail("unexpected recover type")
	}()

	a.Nil(recoverFromError(nil, nil))

	value := js.Global().Get("Error").New()
	value.Set("message", "this is a js error")
	jsErr := js.Error{Value: value}
	err := recoverFromError(jsErr, nil)
	a.IsType(js.Error{}, err)
	a.Contains(err.Error(), "this is a js error")

	var callbackCalled bool
	_ = recoverFromError(jsErr, func() { callbackCalled = true })
	a.True(callbackCalled)

	_ = recoverFromError(errors.New("other error"), nil)
}

func recoverFromError(jsErr any, cb func()) (err error) {
	defer handleJSError(&err, cb)

	if jsErr != nil {
		panic(jsErr)
	}

	return nil
}

func TestWebSocket_SendBinary(t *testing.T) {

	a := assert.New(t)

	mockSocket := js.Global().Get("Object").New()
	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		a.True(len(args) == 1)
		a.True(args[0].Length() == 3)
		a.True(args[0].Index(2).Int() == 2)
		return nil
	})
	defer cb.Release()
	mockSocket.Set("send", cb)

	ws := &webSocket{v: godom.FromJsValue(mockSocket)}
	err := ws.SendBinary([]byte{0, 1, 2})
	a.Nil(err)

}

func TestWebSocket_SendText(t *testing.T) {
	a := assert.New(t)

	mockSocket := js.Global().Get("Object").New()
	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		a.True(len(args) == 1)
		a.Equal("a message", args[0].String())
		return nil
	})
	defer cb.Release()
	mockSocket.Set("send", cb)

	ws := &webSocket{v: godom.FromJsValue(mockSocket)}
	err := ws.SendText("a message")
	a.Nil(err)
}

func TestWebSocket_addEventListener(t *testing.T) {
	a := assert.New(t)
	mockSocket := js.Global().Get("Object").New()

	callMap := make(map[string][]js.Value)

	mockSocket.Set("addEventListener", js.FuncOf(func(this js.Value, args []js.Value) any {
		callMap["addEventListener"] = args
		return nil
	}))
	mockSocket.Set("removeEventListener", js.FuncOf(func(this js.Value, args []js.Value) any {
		callMap["removeEventListener"] = args
		return nil
	}))
	mockSocket.Set("dispatchEvent", js.FuncOf(func(this js.Value, args []js.Value) any {
		callMap["addEventListener"][1].Invoke(js.Global().Get("Event").New("open"))
		return nil
	}))

	ws := &webSocket{v: godom.FromJsValue(mockSocket)}
	var eventCalled bool
	release := ws.addEventListener("open", func(event godom.Value) {
		eventCalled = true
	})
	mockSocket.Call("dispatchEvent", js.Global().Get("Event").New("open"))
	a.True(eventCalled)

	a.Equal(2, len(callMap["addEventListener"]))
	a.Equal(0, len(callMap["removeEventListener"]))

	release()
	a.Equal(2, len(callMap["removeEventListener"]))

}

func Test_message(t *testing.T) {
	a := assert.New(t)

	var expString string
	var expBytes []byte
	w := &webSocket{
		text:   func(message string) { expString = message },
		binary: func(message []byte) { expBytes = message },
	}
	event := js.Global().Get("Object").New()
	event.Set("data", "string data")

	w.message(godom.FromJsValue(event))
	a.Equal("string data", expString)

	event = js.Global().Get("Object").New()
	src := []byte{0, 1, 2, 3, 4}
	uint8Array := js.Global().Get("Uint8Array").New(len(src))
	js.CopyBytesToJS(uint8Array, src)
	event.Set("data", uint8Array)

	w.message(godom.FromJsValue(event))
	a.EqualValues(src, expBytes)

}
