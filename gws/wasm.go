//go:build js && wasm

package gws

import (
	"errors"
	"github.com/mlctrez/godom"
	"syscall/js"
)

func (ws *webSocket) addEventListener(eventType string, fn godom.OnEvent) func() {
	jsFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		fn(godom.FromJsValue(args[0]))
		return nil
	})
	ws.v.Call("addEventListener", eventType, jsFunc)
	return func() {
		ws.v.Call("removeEventListener", eventType, jsFunc)
		jsFunc.Release()
	}
}

func (ws *webSocket) sendText(message string) (err error) {
	defer handleJSError(&err, nil)
	ws.v.Call("send", message)
	return err
}

func (ws *webSocket) sendBinary(message []byte) (err error) {
	defer handleJSError(&err, nil)
	uint8Array := js.Global().Get("Uint8Array").New(len(message))
	js.CopyBytesToJS(uint8Array, message)
	ws.v.Call("send", uint8Array)
	return err
}

func handleJSError(err *error, onErr func()) {
	r := recover()

	switch e := r.(type) {
	case error:
		var jsErr js.Error
		if errors.As(e, &jsErr) {
			*err = jsErr
			if onErr != nil {
				onErr()
			}
			return
		}
	}

	if r != nil {
		panic(r)
	}
}
