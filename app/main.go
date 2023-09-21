package main

import (
	"fmt"
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
	doc := a.Document.DocApi()
	if value == "buttonOne" {
		e.AddEventListener("click", func(event dom.Value) {
			e.SetAttribute("disabled", true)
			go time.AfterFunc(time.Second*1, func() {
				text := fmt.Sprintf("<p>the button was replaced %s</p>", name)
				el := doc.H(text)
				e.ReplaceWith(el)
				go time.AfterFunc(time.Second*2, func() {
					el.ReplaceWith(doc.H(`<button go="buttonOne">click me</button>`))
				})
			})
		})
	}
	if value == "wysiwyg" {
		e.AddEventListener("input", func(event dom.Value) {
			var targetDiv dom.Element
			for _, el := range e.Parent().ChildNodes() {
				if div, ok := el.(dom.Element); ok {
					if div.This().Call("getAttribute", "id").String() == "wysiwyg-target" {
						targetDiv = div.(dom.Element)
					}
				}
			}
			for _, node := range targetDiv.ChildNodes() {
				if el, ok := node.(dom.Element); ok {
					el.Remove()
				}
			}
			targetDiv.AppendChild(doc.H(event.Get("target").Get("value").String()))
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
	case "/edit":
		a.body = doc.H(`<body>` +
			`<a href="/">index page</a>` +
			`<p>Edit page</p>` +
			`<textarea go="wysiwyg" rows="25" cols="80"></textarea>` +
			`<hr/><div id="wysiwyg-target"></div>` +
			`</body>`)
	case "/two":
		a.body = doc.H(`<body>` +
			`<a href="/">index page</a>` +
			`<p>This is page two</p>` +
			`</body>`)
	default:
		a.body = doc.H(`<body>` +
			`<a href="/two">page two</a>` +
			`<br/>` +
			`<a href="/edit">edit page</a>` +
			`<br/>` +
			`<p>This is the index page</p>` +
			`<a href="https://github.com/mlctrez/godom/">outside url</a>` +
			`<br/>` +
			`</body>`)
		a.body.AppendChild(doc.H(`<button go="buttonOne">click me</button>`))

	}
	previousBody.ReplaceWith(a.body)
}
