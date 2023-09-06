package godom

// Window is defined as https://developer.mozilla.org/en-US/docs/Web/API/Window
type Window interface {
	ThisValue
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

type Location interface {
	Reload()
	Href() string
}

var _ Location = (*location)(nil)

type location struct {
	this Value
}

func (l *location) Reload() {
	l.this.Call("reload")
}

func (l *location) Href() string {
	return l.this.Get("href").String()
}

type Console interface {
	Log(args ...any)
	Error(args ...any)
}

var _ Console = (*console)(nil)

type console struct {
	this Value
}

func (c *console) Log(args ...any) {
	if c.this != nil {
		c.this.Call("log", args...)
	}
}

func (c *console) Error(args ...any) {
	if c.this != nil {
		c.this.Call("error", args...)
	}
}
