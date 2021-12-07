package godom

// Window is defined as https://developer.mozilla.org/en-US/docs/Web/API/Window
type Window interface {
	Document() Document
	Navigator() Navigator
	Location() Location
	Console() Console
}

type Navigator interface{}

var _ Navigator = (*navigator)(nil)

type navigator struct {
	this Value
}

type Location interface{}

var _ Location = (*location)(nil)

type location struct {
	this Value
}

type Console interface{}

var _ Console = (*console)(nil)

type console struct {
	this Value
}
