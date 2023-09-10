package main

import (
	"context"
	_ "embed"
	"fmt"
	dom "github.com/mlctrez/godom"
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

	go a.monitorServer()

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
		if _, err := gfet.Fetch(&gfet.Request{URL: href}); err == nil || time.Now().After(endAt) {
			break
		}
		time.Sleep(time.Millisecond * time.Duration(500))
	}
	dom.Global().Get("location").Call("reload")
}

func (a *App) onBinary(message []byte) {
	if string(message) == "wasm" {
		a.ctxCancel()
	}
}

func (a *App) monitorServer() {

	a.ws = gws.New(gws.Rel("ws"))
	a.ws.OnBinaryMessage(a.onBinary)
	a.ws.OnError(dom.EventFunc(a.ctxCancel))
	a.ws.OnClose(gws.CloseFunc(a.ctxCancel))

	pingTicker := time.NewTicker(time.Second)
	defer pingTicker.Stop()
	go func() {
		for {
			select {
			case <-a.ctx.Done():
				return
			case <-pingTicker.C:
				if err := a.ws.SendBinary([]byte("keepalive")); err != nil {
					a.ctxCancel()
					return
				}
			}
		}
	}()

	<-a.ctx.Done()
}

func main() {
	(&App{}).Run()
}
