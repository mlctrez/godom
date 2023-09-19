package base

import (
	"context"
	dom "github.com/mlctrez/godom"
	"github.com/mlctrez/godom/gdutil"
	"github.com/mlctrez/godom/gfet"
	"github.com/mlctrez/godom/gws"
	"time"
)

type App struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	Global    dom.Value
	Window    dom.Value
	events    chan dom.Value
	ws        gws.WebSocket
}

func (a *App) eventHandlers() func() {
	var releases []dom.Func
	toRelease := func(fn dom.Func) dom.Func {
		releases = append(releases, fn)
		return fn
	}
	a.Window.Set("onclick", toRelease(dom.FuncOf(a.click)))
	return func() {
		for _, fn := range releases {
			fn.Release()
		}
	}
}

func (a *App) click(this dom.Value, args []dom.Value) any {
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
		func() {
			defer func() {
				if recover() != nil {
					a.events <- event
				}
			}()
			a.Window.Get("history").Call("pushState", nil, "", target.Get("href"))
			a.events <- a.Window.Get("location")
		}()
		return nil
	}
	a.events <- event
	return nil
}

func (a *App) RunMain(eventHandler func(value dom.Value)) {
	a.ctx, a.ctxCancel = context.WithCancel(context.Background())
	a.events = make(chan dom.Value, 100)
	a.Global = dom.Global()
	a.Window = a.Global.Get("window")

	// TODO: add ability to turn this on and off
	go a.KeepAlive()

	release := a.eventHandlers()
	a.events <- a.Window.Get("location")
	for {
		select {
		case event := <-a.events:
			eventHandler(event)
		case <-a.ctx.Done():
			release()
			a.tryReconnect()
			return
		}
	}
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

	onBinary := func(message []byte) {
		if string(message) == "wasm" {
			a.ctxCancel()
		}
	}

	ws := gws.New(gws.Rel("ws"))
	a.ws = ws
	ws.OnBinaryMessage(onBinary)
	ws.OnError(dom.EventFunc(a.ctxCancel))
	ws.OnClose(gws.CloseFunc(a.ctxCancel))

	gdutil.Periodic(a.ctx, time.Second, func() (ok bool) {
		if err := ws.SendBinary([]byte("keepalive")); err == nil {
			return true
		} else {
			a.ctxCancel()
			return false
		}
	})
}
