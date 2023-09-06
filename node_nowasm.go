//go:build !wasm

package godom

func (d *node) AddEventListener(eventType string, fn OnEvent) func() {
	// TODO: !wasm implementation?
	return func() {
		Global().Console().Log("event listener released")
	}
}
