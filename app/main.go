//go:build js && wasm

package main

import (
	"context"
	_ "embed"
	"fmt"
	dom "github.com/mlctrez/godom"
	"github.com/mlctrez/godom/gws"
	fetch "marwan.io/wasm-fetch"
	"time"
)

type App struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	c         dom.Console
	l         dom.Location
	ws        gws.WebSocket
}

// Run is the main entry point
func (a *App) Run() {
	a.ctx, a.ctxCancel = context.WithCancel(context.Background())
	g := dom.Global()
	a.c = g.Console()
	a.l = g.Location()

	a.c.Log("startup")
	go a.monitorServer()

	document := g.Document()
	doc := dom.Doc{Doc: document}

	var body dom.Node
	body = doc.El("body")
	p := doc.El("p")
	p.AppendChild(doc.Doc.CreateTextNode("click here to close websocket"))
	p.AddEventListener("click", func(event dom.Value) {

		a.ws.Close()
	})

	removedDiv := doc.El("div")
	removedDiv.AppendChild(p)
	body.AppendChild(removedDiv)

	buttonOne := doc.H(`<button>click to remove div</button>`)
	buttonOne.AddEventListener("click", func(event dom.Value) {
		a.c.Log("button one", event)
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
		if _, err := fetch.Fetch(a.l.Href(), &fetch.Opts{}); err == nil || time.Now().After(endAt) {
			break
		}
		time.Sleep(time.Millisecond * time.Duration(500))
	}
	a.l.Reload()
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
					fmt.Println("error on send", err)
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
