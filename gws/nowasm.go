//go:build !js && !wasm

package gws

import "github.com/mlctrez/godom"

func (ws *webSocket) addEventListener(eventType string, fn godom.OnEvent) func() {
	return func() {}
}

func (ws *webSocket) sendBinary(message []byte) (err error) {
	return nil
}

func (ws *webSocket) sendText(message string) (err error) {
	return nil
}
