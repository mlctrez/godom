package api

import (
	"github.com/mlctrez/godom"
	"net/url"
)

type Context struct {
	Doc    godom.Doc
	URL    *url.URL
	Events chan Event
}

type Handler interface {
	Headers(ctx *Context, header godom.Element)
	Body(ctx *Context) godom.Element
}

type Event interface {
}

type Location struct {
	Event
	URL      *url.URL
	External bool
	PopState bool
}
