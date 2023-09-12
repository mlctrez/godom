package main

import (
	"context"
	_ "embed"
	dom "github.com/mlctrez/godom"
	"github.com/mlctrez/godom/gdutil"
	"github.com/mlctrez/godom/gfet"
	"github.com/mlctrez/godom/gws"
	"strings"
	"time"
)

type App struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	ws        gws.WebSocket
}

//go:embed body.html
var bodyString string

type AppEvent struct {
	el     dom.Element
	event  string
	dataGo string
}

func (e *AppEvent) handleEvent(event dom.Value) {
	dom.Console().Log("%s %o %o", e.dataGo, event, e.el)
}

// Run is the main entry point
func (a *App) Run() {
	a.ctx, a.ctxCancel = context.WithCancel(context.Background())

	go a.KeepAlive()
	c := dom.Console()

	document := dom.Document()
	doc := dom.Document().DocApi()
	doc.CallBack = func(e dom.Element, dataGo []string) {
		for _, val := range dataGo {
			split := strings.Split(val, ".")
			if len(split) != 2 {
				c.Log("invalid data-go attribute %s on %o", val, e.This())
				continue
			}
			ae := &AppEvent{el: e, event: split[1], dataGo: split[0]}
			e.AddEventListener(split[1], ae.handleEvent)
		}
	}

	var body dom.Element
	body = doc.H(bodyString)

	document.Body().ReplaceWith(body)

	<-a.ctx.Done()
	a.tryReconnect()
}

func (a *App) tryReconnect() {
	endAt := time.Now().Add(time.Second * 5)
	for {
		href := dom.Global().Get("location").Get("href").String()
		req := &gfet.Request{URL: href, Method: "OPTIONS"}
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
			ok = true
		} else {
			a.ctxCancel()
		}
		return ok
	})
}

func main() {
	(&App{}).Run()
}
