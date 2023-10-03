package wasm

import (
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/nap"
	"net/url"
	"strings"
)

type App struct {
	c        *nap.Config
	Window   godom.Value
	events   chan nap.Event
	releases []godom.Func
}

func Run(c *nap.Config) (err error) {
	return (&App{c: c}).run()
}

func (a *App) setup() (err error) {
	a.c.Logger.Debug("setup")
	a.events = make(chan nap.Event, 100)
	a.Window = godom.Global().Get("window")
	a.windowEventHandlers()
	return nil
}

func (a *App) parseUrl(href string) (u *url.URL) {
	var err error
	if u, err = url.Parse(href); err != nil {
		a.c.Logger.Error("setup", "error", err, "href", href)
		return &url.URL{Path: "/"}
	}
	a.c.Logger.Debug("parseUrl", "href", href, "path", u.Path)
	return u
}

func (a *App) run() (err error) {
	if err = a.setup(); err != nil {
		return err
	}

	a.c.App.Init(a.c, godom.Document().DocumentElement())
	u := a.URL()
	u.Path = strings.TrimPrefix(u.Path, a.c.PagesPath)
	u.Path = strings.TrimSuffix(u.Path, ".html")
	a.events <- &nap.Location{URL: u, PopState: true}
	for {
		select {
		case value := <-a.events:
			switch v := value.(type) {
			case *nap.Location:
				a.location(v)
			}
		case <-a.c.Context.Done():
			a.releaseEventHandlers()
			//if devReloadEnabled {
			//	a.tryReconnect()
			//}
			return a.c.Context.Err()
		}
	}
}

func (a *App) location(v *nap.Location) {

	if !v.External && !v.PopState {
		a.Window.Get("history").Call("pushState", nil, "", a.historyUrl(v.URL))
	}
	if !v.External {
		a.navigate(v.URL)
	} else {
		// TODO: external url handling
	}

}

func (a *App) navigate(u *url.URL) {
	a.c.Logger.Debug("navigate", "url", u.String())
	ctx := nap.NewDocContext(u, a.events)
	for _, page := range a.c.Pages {
		if page.Regexp.MatchString(u.Path) {
			a.replaceBody(page.PageFunc(ctx))
			return
		}
	}
	if a.c.NotFoundFunc != nil {
		a.replaceBody(a.c.NotFoundFunc(u, ctx))
		return
	}
	a.replaceBody(ctx.DocApi().H("<body>page not found</body>"))
}

func (a *App) replaceBody(htmlPage godom.Element) {
	a.c.Logger.Debug("replaceBody", "page", htmlPage.String())
	body := htmlPage.GetElementsByTagName("body")[0]
	htmlPage.RemoveChild(body.This())

	body.SetParent(godom.Document().Body().Parent())
	godom.Document().Body().ReplaceWith(body)
}

func (a *App) historyUrl(u *url.URL) string {
	histUrl := &url.URL{}
	*histUrl = *u
	if histUrl.Path == "/" {
		histUrl.Path = "/index"
	}
	histUrl.Path = a.c.PagesPath + histUrl.Path + ".html"
	return histUrl.String() + "?_ij_reload=RELOAD_ON_SAVE"
}

func (a *App) URL() *url.URL {
	return a.parseUrl(a.Window.Get("location").Get("href").String())
}

func (a *App) addRelease(fn godom.Func) godom.Func {
	a.releases = append(a.releases, fn)
	return fn
}

func (a *App) releaseEventHandlers() {
	for _, release := range a.releases {
		release.Release()
	}
}

func (a *App) windowEventHandlers() {
	a.c.Logger.Debug("windowEventHandlers")
	a.Window.Set("onclick", a.addRelease(godom.FuncOf(a.onClick)))
	a.Window.Set("onpopstate", a.addRelease(godom.FuncOf(a.onPopState)))
	// add additional handlers here
}

func (a *App) onPopState(this godom.Value, args []godom.Value) any {
	a.c.Logger.Debug("onPopState")
	u := a.parseUrl(a.Window.Get("location").Get("href").String())
	u.Path = strings.TrimPrefix(u.Path, a.c.PagesPath)
	u.Path = strings.TrimSuffix(u.Path, ".html")
	a.events <- &nap.Location{URL: u, PopState: true}
	return nil
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
		next := a.parseUrl(href)
		if next.Host != a.URL().Host {
			a.events <- &nap.Location{URL: next, External: true}
		} else {
			next.Path = strings.TrimPrefix(next.Path, a.c.PagesPath)
			a.events <- &nap.Location{URL: next}
		}
		return nil
	}
	return nil
}
