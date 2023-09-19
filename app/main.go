package main

import (
	_ "embed"
	"fmt"
	dom "github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app/base"
	"net/url"
)

type App struct {
	base.App
	body dom.Element
}

func main() {
	a := &App{}
	a.RunMain(a.handleEvent)
}

func (a *App) handleEvent(value dom.Value) {
	console := a.Window.Get("console")
	//console.Call("log", "handleEvent", value)
	if value.InstanceOf(a.Global.Get("Location")) {
		u, err := url.Parse(value.Get("href").String())
		if err != nil {
			console.Call("error", err.Error())
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

func (a *App) navigate(u *url.URL) {
	//fmt.Println(u)

	document := dom.Document()
	doc := dom.Document().DocApi()
	doc.CallBack = func(e dom.Element, dataGo string) {
		e.AddEventListener("click", func(event dom.Value) {
			fmt.Println(e, dataGo)
		})
	}
	previousBody := a.body
	if previousBody == nil {
		previousBody = document.Body()
	}

	switch u.Path {
	case "/":
		a.body = doc.H(`<body>` +
			`<a href="/two">page two</a>` +
			`<p>This is the index page</p>` +
			`<a href="https://github.com/mlctrez/godom/">outside url</a>` +
			`<br/><button go="buttonOne">click me</button>` +
			`</body>`)

	case "/two":
		a.body = doc.H(`<body>` +
			`<a href="/">index page</a>` +
			`<p>This is page two</p>` +
			`</body>`)
	}
	previousBody.ReplaceWith(a.body)

}
