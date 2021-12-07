package godom

type EventTarget interface {
	AddEventListener(eventName string, listener interface{})
}

var _ EventTarget = (*eventTarget)(nil)

type eventTarget struct{}

func (d *eventTarget) AddEventListener(eventName string, listener interface{}) {
	panic(IM)
}
