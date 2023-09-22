package wsm

import (
	"context"
	dom "github.com/mlctrez/godom"
	"github.com/mlctrez/godom/gdutil"
	"github.com/mlctrez/godom/gfet"
	"github.com/mlctrez/godom/gsrv/api"
	"github.com/mlctrez/godom/gws"
	"net/url"
	"time"
)

func Run(h api.Handler) {
	a := &App{}
	a.ctx, a.ctxCancel = context.WithCancel(context.Background())
	a.h = h
	a.events = make(chan api.Event, 100)
	a.Global = dom.Global()
	a.Window = a.Global.Get("window")
	a.Document = dom.Document()
	a.lastBody = a.Document.Body()

	// TODO: add ability to turn this on and off
	go a.KeepAlive()

	release := a.eventHandlers()
	u, err := url.Parse(a.Window.Get("location").Get("href").String())
	if err != nil {
		a.Window.Get("console").Call("error", "location.href parse error", err.Error())
		u = &url.URL{Path: "/"}
	}

	a.h.Headers(&api.Context{Doc: a.Document.DocApi(), URL: u}, a.Document.Head())

	a.events <- &api.Location{URL: u}
	for {
		select {
		case value := <-a.events:
			switch v := value.(type) {
			case *api.Location:
				if !v.External && !v.PopState {
					a.Window.Get("history").Call("pushState", nil, "", v.URL.String())
				}
				if !v.External {
					newBody := a.h.Body(&api.Context{Doc: a.Document.DocApi(), URL: v.URL, Events: a.events})
					a.lastBody.ReplaceWith(newBody)
					a.lastBody = newBody
				}
			}
		case <-a.ctx.Done():
			release()
			a.tryReconnect()
			return
		}
	}
}

type App struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	h         api.Handler
	Global    dom.Value
	Window    dom.Value
	Document  dom.DocumentApi
	events    chan api.Event
	ws        gws.WebSocket
	lastBody  dom.Element
}

func (a *App) eventHandlers() func() {
	var releases []dom.Func
	toRelease := func(fn dom.Func) dom.Func {
		releases = append(releases, fn)
		return fn
	}

	a.Window.Set("onclick", toRelease(dom.FuncOf(a.onClick)))
	a.Window.Set("onpopstate", toRelease(dom.FuncOf(a.onPopState)))
	// add additional handlers here

	return func() {
		for _, fn := range releases {
			fn.Release()
		}
	}
}

func (a *App) onClick(this dom.Value, args []dom.Value) any {
	event := args[0]
	target := event.Get("target")
	if !target.Truthy() {
		a.Window.Get("console").Call("error", "target", event)
		return nil
	}
	nodeName := target.Get("nodeName")
	if !nodeName.Truthy() {
		a.Window.Get("console").Call("error", "nodeName", event)
		return nil
	}
	if nodeName.String() == "A" {
		event.Call("preventDefault")
		href := target.Get("href").String()
		u, err := url.Parse(href)
		if err != nil {
			a.Window.Get("console").Call("error", "target.href", err)
			return nil
		}
		wu, _ := url.Parse(a.Window.Get("location").Get("href").String())
		if u.Host != wu.Host {
			a.events <- &api.Location{URL: u, External: true}
		} else {
			a.events <- &api.Location{URL: u}
		}
		return nil
	}
	return nil
}

func (a *App) onPopState(this dom.Value, args []dom.Value) any {
	u, _ := url.Parse(a.Window.Get("location").Get("href").String())
	a.events <- &api.Location{URL: u, PopState: true}
	return nil
}

func (a *App) tryReconnect() {
	endAt := time.Now().Add(time.Second * 5)
	for {
		req := &gfet.Request{URL: "/", Method: "HEAD"}
		_, err := req.Fetch()
		if err == nil || time.Now().After(endAt) {
			break
		}
		time.Sleep(time.Millisecond * time.Duration(500))
	}
	dom.Global().Get("location").Call("reload")
}

func (a *App) KeepAlive() {

	a.ws = gws.New(gws.Rel("ws"))
	defer a.ws.Close()

	a.ws.OnBinaryMessage(func(message []byte) {
		if string(message) == "wasm" {
			a.ctxCancel()
		}
	})
	a.ws.OnError(dom.EventFunc(a.ctxCancel))
	a.ws.OnClose(gws.CloseFunc(a.ctxCancel))

	gdutil.Periodic(a.ctx, time.Second, func() (ok bool) {
		if err := a.ws.SendBinary([]byte("keepalive")); err == nil {
			return true
		} else {
			a.ctxCancel()
			return false
		}
	})

}
