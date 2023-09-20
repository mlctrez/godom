package main

import (
	dom "github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app/base"
	"net/url"
	"time"
)

type App struct {
	base.App
	body dom.Element
}

func main() {
	a := &App{}
	a.RunMain(a.eventHandler)
}

func (a *App) eventHandler(value dom.Value) {
	console := a.Window.Get("console")
	//console.Call("log", "eventHandler", value)
	if value.InstanceOf(a.Global.Get("Location")) {
		u, err := url.Parse(value.Get("href").String())
		if err != nil {
			console.Call("error", err.Error())
			return
		}
		a.navigate(u)
	}
	if value.InstanceOf(a.Global.Get("PointerEvent")) {
		if value.Get("target").Get("nodeName").String() == "A" {
			// for external urls
			console.Call("log", value.Get("target").Get("href"))
		}
	}
}

func (a *App) docCallback(e dom.Element, name, value string) {
	if value == "buttonOne" {
		e.AddEventListener("click", func(event dom.Value) {
			e.SetAttribute("disabled", true)
			go time.AfterFunc(time.Second*1, func() {
				e.ReplaceWith(a.Document.DocApi().T("the button was replaced " + name))
			})
		})
	}
}

func (a *App) navigate(u *url.URL) {
	doc := dom.Document().DocApi()
	doc.CallBack = a.docCallback
	previousBody := a.body
	if previousBody == nil {
		previousBody = a.Document.Body()
	}

	switch u.Path {
	case "/two":
		a.body = doc.H(`<body>` +
			`<a href="/">index page</a>` +
			`<p>This is page two</p>` +
			`</body>`)
	default:
		a.body = doc.H(`<body>` +
			`<a href="/two">page two</a>` +
			`<p>This is the index page</p>` +
			`<a href="https://github.com/mlctrez/godom/">outside url</a>` +
			`<br/>` +
			`</body>`)
		a.body.AppendChild(doc.H(`<button go="buttonOne">click me</button>`))

	}
	previousBody.ReplaceWith(a.body)
}
