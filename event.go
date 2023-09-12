package godom

type Event interface {
}

type EventTarget interface {
	AddEventListener(eventType string, fn OnEvent) func()
	RemoveEventListener(eventType string, fn OnEvent) func()
	DispatchEvent(event Event)
}

type EventListener interface {
	AddEventListener(eventType string, fn OnEvent) func()
}

type OnEvent func(event Value)

func EventFunc(fn func()) OnEvent {
	return func(event Value) { fn() }
}
