package main

import (
	"context"
	_ "embed"
	"fmt"
	dom "github.com/mlctrez/godom"
	"github.com/mlctrez/godom/gdutil"
	"github.com/mlctrez/godom/gfet"
	"github.com/mlctrez/godom/gws"
	"time"
)

type App struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	ws        gws.WebSocket
}

// Run is the main entry point
func (a *App) Run() {
	a.ctx, a.ctxCancel = context.WithCancel(context.Background())

	go a.KeepAlive()

	document := dom.Document()
	doc := dom.Document().DocApi()

	var body dom.Node
	body = doc.El("body")
	p := doc.El("p", doc.At("style", "cursor:pointer"))
	p.AppendChild(document.CreateTextNode("click here to close websocket"))
	p.AddEventListener("click", func(event dom.Value) {
		a.ws.Close()
	})

	removedDiv := doc.El("div")
	removedDiv.AppendChild(p)
	body.AppendChild(removedDiv)

	buttonOne := doc.H(`<button>click to remove div</button>`)
	buttonOne.AddEventListener("click", func(event dom.Value) {
		fmt.Println("button one")
		removedDiv.Remove()
	})

	body.AppendChild(buttonOne)

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

	go gdutil.Periodic(a.ctx, time.Second, func() (ok bool) {
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
