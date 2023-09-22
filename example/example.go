package example

import (
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app"
	"time"
)

var _ app.Handler = (*example)(nil)

type example struct {
}

func New() app.Handler {
	return &example{}
}

func (e *example) Prepare(ctx *app.ServerContext) {
	ctx.Main = "example/bin/main.go"
	ctx.Output = "build/app.wasm"
	ctx.Watch = []string{"example", "app"}
	ctx.Address = ":8080"
}

var style = `<style>
body { color: white; background-color: black; } 
</style>`

func (e *example) Headers(ctx *app.Context, header godom.Element) {
	for _, node := range header.ChildNodes() {
		if node.NodeName() == "title" {
			node.This().Set("innerHTML", "godom")
		}
	}
	if len(header.GetElementsByTagName("style")) == 0 {
		header.AppendChild(ctx.Doc.H(style))
	}
}

func (e *example) Body(ctx *app.Context) godom.Element {
	switch ctx.URL.Path {
	case "/":
		return e.index(ctx)
	default:
		return ctx.Doc.H(pageNotFound)
	}
}

var indexPage = `<body>
<table>
<tr>
  <td><button go="button">example one</button><div go="list"/></td>
</tr>
</table>
</body>`

type exampleOne struct {
	button godom.Element
	clear  godom.Element
	list   godom.Element
}

func (eo *exampleOne) mapper(e godom.Element, name, data string) {
	if name != "go" {
		return
	}
	switch data {
	case "button":
		eo.button = e
	case "clear":
		eo.clear = e
	case "list":
		eo.list = e
	}
}

var exOneRow = `<tr><td>
<button go="button">example one</button>
<button go="clear">clear</button> 
<div go="list"/>
</td></tr>`

func (eo *exampleOne) render(ctx *app.Context) godom.Element {
	ctx.Doc.CallBack = eo.mapper
	row := ctx.Doc.H(exOneRow)
	eo.button.AddEventListener("click", func(event godom.Value) {
		list := eo.list
		list.AppendChild(ctx.Doc.El("br"))
		list.AppendChild(ctx.Doc.H(`<span>` + time.Now().Format(time.RFC3339Nano) + `</span>`))
		if len(list.ChildNodes()) > 12 {
			list.RemoveChild(list.ChildNodes()[0].This())
			list.RemoveChild(list.ChildNodes()[0].This())
		}
	})
	eo.clear.AddEventListener("click", func(event godom.Value) {
		for len(eo.list.ChildNodes()) > 0 {
			eo.list.RemoveChild(eo.list.ChildNodes()[0].This())
		}
	})
	return row
}

func (e *example) index(ctx *app.Context) godom.Element {
	body := ctx.Doc.H("<body><table></table></body>")
	body.ChildNodes()[0].AppendChild((&exampleOne{}).render(ctx))
	return body
}

var pageNotFound = `<body>
<p style="color:red">page not found</p>
<br/>
Go back to <a href="/">index page</a>
</body>`
