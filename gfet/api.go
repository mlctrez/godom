package gfet

import (
	"context"
	_ "embed"
	"errors"
	"github.com/mlctrez/godom"
	"net/textproto"
)

type fetch struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	res        *Response
	err        error
	releases   []godom.Func
}

type Request struct {
	URL            string
	Method         string
	Headers        textproto.MIMEHeader
	Mode           string
	Credentials    string
	Cache          string
	Referrer       string
	ReferrerPolicy string
	Integrity      string
	KeepAlive      *bool
	Body           []byte
}

type Response struct {
	Ok         bool
	Redirected bool
	Status     int
	StatusText string
	Type       string
	URL        string
	BodyUsed   bool
	Headers    textproto.MIMEHeader
	Body       []byte
}

func (f *fetch) fulfilled(arg godom.Value) {
	f.res = &Response{Headers: textproto.MIMEHeader{}}

	headersIt := arg.Get("headers").Call("entries")
	for {
		n := headersIt.Call("next")
		if n.Get("done").Bool() {
			break
		}
		pair := n.Get("value")
		key, value := pair.Index(0).String(), pair.Index(1).String()
		f.res.Headers.Add(key, value)
	}

	f.res.Ok = arg.Get("ok").Bool()
	f.res.Redirected = arg.Get("redirected").Bool()
	f.res.Status = arg.Get("status").Int()
	f.res.StatusText = arg.Get("statusText").String()
	f.res.Type = arg.Get("type").String()
	f.res.URL = arg.Get("url").String()
	f.res.BodyUsed = arg.Get("bodyUsed").Bool()

	arg.Call("arrayBuffer").Call("then", f.funcOf(f.arrayBuffer))
}

func (f *fetch) arrayBuffer(arg godom.Value) {
	defer f.cancelFunc()
	f.res.Body = arg.Bytes()
}

func (f *fetch) rejected(arg godom.Value) {
	defer f.cancelFunc()
	f.err = errors.New(arg.Get("message").String())
}

func (f *fetch) funcOf(target func(arg godom.Value)) godom.Func {
	fn := godom.FuncOf(func(this godom.Value, args []godom.Value) any {
		target(args[0])
		return nil
	})
	f.releases = append(f.releases, fn)
	return fn
}

func (f *fetch) release() {
	for _, fn := range f.releases {
		fn.Release()
	}
}

func Fetch(r *Request) (res *Response, err error) {
	f := &fetch{}
	f.ctx, f.cancelFunc = context.WithCancel(context.TODO())
	defer f.release()

	// TODO: use context with timeout and cancel fetch

	global := godom.Global()
	optionsMap := global.Get("Object").New()
	if r.Headers != nil {
		headersMap := global.Get("Object").New()
		for key := range r.Headers {
			headersMap.Set(key, r.Headers.Get(key))
		}
		optionsMap.Set("headers", headersMap)
	}
	// TODO: other request options

	fetchApi := global.Get("fetch")
	fetchApi.Invoke(r.URL, optionsMap).
		Call("then", f.funcOf(f.fulfilled)).
		Call("catch", f.funcOf(f.rejected))

	<-f.ctx.Done()
	return f.res, f.err
}
