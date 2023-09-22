package example

import (
	"fmt"
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
	for i, node := range header.ChildNodes() {
		if node.NodeName() == "title" {
			header.ChildNodes()[i] = ctx.Doc.H("<title>godom</title>")
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

type exampleOne struct {
	button godom.Element
	reset  godom.Element
	list   godom.Element
}

func (e *example) index(ctx *app.Context) godom.Element {
	body := ctx.Doc.H("<body><table><tbody/></table></body>")
	body.GetElementsByTagName("tbody")[0].AppendChild((&exampleOne{}).render(ctx))
	return body
}

func (eo *exampleOne) mapper(e godom.Element, name, data string) {
	if name != "go" {
		return
	}
	// TODO: figure out how to do json style bindings for this
	switch data {
	case "button":
		eo.button = e
	case "reset":
		eo.reset = e
	case "list":
		eo.list = e
	}
}

var exOneRow = `<tr><td>
<button go="button">example one</button><button go="reset">reset</button> 
<div go="list"></div>
</td></tr>`

func (eo *exampleOne) render(ctx *app.Context) godom.Element {

	doc := godom.Doc{Doc: ctx.Doc.Doc, CallBack: eo.mapper}
	row := doc.H(exOneRow)
	list := eo.list
	eo.button.AddEventListener("click", func(event godom.Value) {
		eo.button.SetAttribute("disabled", true)
		go time.AfterFunc(1*time.Second, func() {
			eo.button.RemoveAttribute("disabled")
		})
		list.AppendChild(doc.El("br"))
		span := fmt.Sprintf("<span>%s</span>", time.Now().Format(time.RFC3339Nano))
		list.AppendChild(doc.H(span))
		if len(list.ChildNodes()) > 12 {
			list.RemoveChild(list.ChildNodes()[0].This())
			list.RemoveChild(list.ChildNodes()[0].This())
		}
	})
	eo.reset.AddEventListener("click", func(event godom.Value) {
		for len(list.ChildNodes()) > 0 {
			list.RemoveChild(eo.list.ChildNodes()[0].This())
		}
	})
	return row
}

var pageNotFound = `<body>
<p style="color:red">page not found</p>
<br/>
Go back to <a href="/">index page</a>
</body>`
