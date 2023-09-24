package app

import (
	"github.com/mlctrez/godom"
	"io"
	"net/url"
	"runtime"
)

type Context struct {
	Doc    godom.DocApi
	URL    *url.URL
	Events chan Event
}

func (c *Context) IsWasm() bool {
	return runtime.GOARCH == "wasm"
}

type ServerContext struct {
	Main         string
	Output       string
	Address      string
	Watch        []string
	ShowWasmSize bool
}

type Handler interface {
	Prepare(ctx *ServerContext)
	Headers(ctx *Context, header godom.Element)
	Body(ctx *Context) godom.Element

	Serve(request Request, response Response) bool
}

type Request interface {
	Method() string
	URL() *url.URL
	Headers() map[string]string
	Body() io.ReadCloser
}

type Response interface {
	SetHeader(name, value string)
	WriteHeader(statusCode int)
	Write([]byte) (int, error)
}

type Event interface {
}

type Location struct {
	Event
	URL      *url.URL
	External bool
	PopState bool
}
