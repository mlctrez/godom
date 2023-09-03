package godom

type Event interface {
}

type EventListener interface {
	AddEventListener(eventType string, fn OnEvent) func()
}

type OnEvent func(event Value) any
