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
	ws        gws.WebSocket
}

func (a *App) RunMain(mainApp func()) {
	a.ctx, a.ctxCancel = context.WithCancel(context.Background())
	go a.KeepAlive()
	mainApp()
	<-a.ctx.Done()
	a.tryReconnect()
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
