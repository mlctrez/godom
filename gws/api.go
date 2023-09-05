package gws

import "github.com/mlctrez/godom"

// CloseEvent is the type passed to a WebSocket close handler.
type CloseEvent struct {
	Code     uint16
	Reason   string
	WasClean bool
}

type WebSocket interface {
	OnOpen(fn godom.OnEvent)
	OnError(fn godom.OnEvent)

	OnBinaryMessage(func(message []byte))
	OnTextMessage(func(message string))

	SendBinary(message []byte) error
	SendText(message string) error

	OnClose(fn func(event CloseEvent))
	Close()
}

func CloseFunc(fn func()) func(event CloseEvent) {
	return func(event CloseEvent) { fn() }
}
