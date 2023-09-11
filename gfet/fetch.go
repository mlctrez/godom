package gfet

import (
	"context"
	_ "embed"
	"errors"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/std"
	"io"
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
	// URL is the first argument passed to fetch.
	URL string

	// Method is the http verb (constants are copied from net/http to avoid import)
	Method string

	// Headers is a map of http headers to send.
	Headers map[string]string

	// Body is the body request
	Body io.Reader

	// Mode docs https://developer.mozilla.org/en-US/docs/Web/API/Request/mode
	Mode string

	// Credentials docs https://developer.mozilla.org/en-US/docs/Web/API/Request/credentials
	Credentials string

	// Cache docs https://developer.mozilla.org/en-US/docs/Web/API/Request/cache
	Cache string

	// Redirect docs https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch
	Redirect string

	// Referrer docs https://developer.mozilla.org/en-US/docs/Web/API/Request/referrer
	Referrer string

	// ReferrerPolicy docs https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch
	ReferrerPolicy string

	// Integrity docs https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity
	Integrity string

	// KeepAlive docs https://developer.mozilla.org/en-US/docs/Web/API/WindowOrWorkerGlobalScope/fetch
	KeepAlive *bool

	// Signal docs https://developer.mozilla.org/en-US/docs/Web/API/AbortSignal
	Signal context.Context
}

func mapHeaders(mp map[string]string) map[string]interface{} {
	newMap := map[string]interface{}{}
	for k, v := range mp {
		newMap[k] = v
	}
	return newMap
}

func optionsMap(r *Request) (map[string]interface{}, error) {
	mp := map[string]interface{}{}

	if r.Method != "" {
		mp["method"] = r.Method
	}
	if r.Headers != nil {
		mp["headers"] = mapHeaders(r.Headers)
	}
	if r.Mode != "" {
		mp["mode"] = r.Mode
	}
	if r.Credentials != "" {
		mp["credentials"] = r.Credentials
	}
	if r.Cache != "" {
		mp["cache"] = r.Cache
	}
	if r.Redirect != "" {
		mp["redirect"] = r.Redirect
	}
	if r.Referrer != "" {
		mp["referrer"] = r.Referrer
	}
	if r.ReferrerPolicy != "" {
		mp["referrerPolicy"] = r.ReferrerPolicy
	}
	if r.Integrity != "" {
		mp["integrity"] = r.Integrity
	}
	if r.KeepAlive != nil {
		mp["keepalive"] = *r.KeepAlive
	}

	if r.Body != nil {
		bts, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		mp["body"] = string(bts)
	}

	return mp, nil
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
	r := &Response{Headers: textproto.MIMEHeader{}}

	headersIt := arg.Get("headers").Call("entries")
	std.MapEach(headersIt, func(key, val godom.Value) {
		r.Headers.Add(key.String(), val.String())
	})

	r.Ok = arg.Get("ok").Bool()
	r.Redirected = arg.Get("redirected").Bool()
	r.Status = arg.Get("status").Int()
	r.StatusText = arg.Get("statusText").String()
	r.Type = arg.Get("type").String()
	r.URL = arg.Get("url").String()
	r.BodyUsed = arg.Get("bodyUsed").Bool()

	f.res = r
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

func (r *Request) Fetch() (res *Response, err error) {

	var opts map[string]interface{}
	if opts, err = optionsMap(r); err != nil {
		return nil, err
	}

	f := &fetch{}
	f.ctx, f.cancelFunc = context.WithCancel(context.TODO())
	defer f.release()

	if r.Signal != nil {
		controller := godom.Global().Get("AbortController").New()
		signal := controller.Get("signal")
		opts["signal"] = godom.ToJsValue(signal)
		go func() {
			select {
			case <-r.Signal.Done():
				controller.Call("abort")
				return
			case <-f.ctx.Done():
				return
			}
		}()
	}

	fetchApi := godom.Global().Get("fetch")
	go fetchApi.Invoke(r.URL, opts).
		Call("then", f.funcOf(f.fulfilled)).
		Call("catch", f.funcOf(f.rejected))

	<-f.ctx.Done()
	return f.res, f.err
}
