package wsm

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app"
	"github.com/mlctrez/godom/callback"
	"github.com/mlctrez/godom/gdutil"
	"github.com/mlctrez/godom/gfet"
	"github.com/mlctrez/godom/gws"
)

func Run(h app.Handler) {
	a := &App{}
	a.ctx, a.ctxCancel = context.WithCancel(context.Background())
	a.h = h
	a.events = make(chan app.Event, 100)
	a.Global = godom.Global()
	a.Window = a.Global.Get("window")

	href := a.Window.Get("location").Get("href").String()

	devReloadEnabled := strings.Contains(href, "localhost")
	if devReloadEnabled {
		go a.KeepAlive()
	}

	release := a.eventHandlers()

	u, err := url.Parse(href)
	if err != nil {
		a.Window.Get("console").Call("error", "location.href parse error", err.Error())
		u = &url.URL{Path: "/"}
	}

	a.Document = godom.Document()
	a.lastBody = a.Document.Body()

	// TODO: need a way to re-bind elements already in the dom
	callback.Reflect(a.h)(a.Document.DocumentElement(), "go", "html")

	a.prefix = gdutil.GetPrefix(u)
	u.Path = strings.TrimPrefix(u.Path, a.prefix)
	a.events <- &app.Location{URL: u, PopState: true}
	for {
		select {
		case value := <-a.events:
			switch v := value.(type) {
			case *app.Location:
				if !v.External && !v.PopState {
					a.Window.Get("history").Call("pushState", nil, "", a.prefixUrl(v.URL))
				}
				if !v.External {
					api := a.Document.DocApi().WithCallback(callback.Reflect(a.h))
					bCtx := &app.Context{Doc: api, URL: v.URL, Events: a.events}
					newBody := a.h.Body(bCtx)
					a.lastBody.ReplaceWith(newBody)
					a.lastBody = newBody
				} else {
					// TODO: external url handling
				}
			}
		case <-a.ctx.Done():
			release()
			if devReloadEnabled {
				a.tryReconnect()
			}
			return
		}
	}
}

type App struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	h         app.Handler
	Global    godom.Value
	Window    godom.Value
	Document  godom.DocumentApi
	events    chan app.Event
	ws        gws.WebSocket
	lastBody  godom.Element
	prefix    string
}

func (a *App) eventHandlers() func() {
	var releases []godom.Func
	toRelease := func(fn godom.Func) godom.Func {
		releases = append(releases, fn)
		return fn
	}

	a.Window.Set("onclick", toRelease(godom.FuncOf(a.onClick)))
	a.Window.Set("onpopstate", toRelease(godom.FuncOf(a.onPopState)))
	// add additional handlers here

	return func() {
		for _, fn := range releases {
			fn.Release()
		}
	}
}

func (a *App) prefixUrl(u *url.URL) string {
	if a.prefix == "" {
		return u.String()
	}
	// avoid modifying original
	withPrefix := &url.URL{}
	*withPrefix = *u
	withPrefix.Path = a.prefix + withPrefix.Path
	return withPrefix.String()
}

func (a *App) onClick(this godom.Value, args []godom.Value) any {
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
		href := target.Get("href").String()
		if href == "" {
			return nil
		}
		event.Call("preventDefault")
		u, err := url.Parse(href)
		if err != nil {
			a.Window.Get("console").Call("error", "target.href", err)
			return nil
		}
		//fmt.Println("click A", href, u.String())
		wu, _ := url.Parse(a.Window.Get("location").Get("href").String())
		if u.Host != wu.Host {
			a.events <- &app.Location{URL: u, External: true}
		} else {
			a.events <- &app.Location{URL: u}
		}
		return nil
	}
	return nil
}

func (a *App) onPopState(this godom.Value, args []godom.Value) any {
	u, _ := url.Parse(a.Window.Get("location").Get("href").String())
	a.events <- &app.Location{URL: u, PopState: true}
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
	godom.Global().Get("location").Call("reload")
}

func (a *App) KeepAlive() {

	a.ws = gws.New(gws.Rel("ws"))
	defer a.ws.Close()

	a.ws.OnBinaryMessage(func(message []byte) {
		if string(message) == "wasm" {
			a.ctxCancel()
		}
	})
	a.ws.OnError(godom.EventFunc(a.ctxCancel))
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
