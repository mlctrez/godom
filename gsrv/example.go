package gsrv

import (
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/gsrv/api"
	"time"
)

var _ api.Handler = (*Example)(nil)

type Example struct {
}

func (e *Example) Headers(ctx *api.Context, header godom.Element) {
	var hasStyle bool
	for _, node := range header.ChildNodes() {
		if node.NodeName() == "title" {
			node.This().Set("innerHTML", "godom")
		}
		if node.NodeName() == "style" {
			hasStyle = true
		}
	}
	if !hasStyle {
		header.AppendChild(ctx.Doc.H(`<style> body { color: white; background-color: black; } </style>`))
	}
}

func (e *Example) Body(ctx *api.Context) godom.Element {
	switch ctx.URL.Path {
	case "/":
		return e.index(ctx)
	default:
		return ctx.Doc.H(pageNotFound)
	}
}

var indexPage = `<body>
<button go="button">add item</button>
<ul go="list"/>
</body>`

func (e *Example) index(ctx *api.Context) godom.Element {

	em := make(map[string]godom.Element)

	ctx.Doc.CallBack = func(e godom.Element, name, data string) { em[data] = e }
	page := ctx.Doc.H(indexPage)
	if len(em) == 2 {
		em["button"].AddEventListener("click", func(event godom.Value) {
			list := em["list"]

			list.AppendChild(ctx.Doc.H(`<li>` + time.Now().Format(time.RFC3339Nano) + `</li>`))
			if len(list.ChildNodes()) > 5 {
				list.RemoveChild(list.ChildNodes()[0].This())
			}
		})
	}

	return page
}

var pageNotFound = `<body>
<p style="color:red">page not found</p>
<br/>
Go back to <a href="/">index page</a>
</body>`
