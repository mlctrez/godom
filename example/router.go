package example

import (
	_ "embed"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app"
	"github.com/mlctrez/godom/example/demos"
	"github.com/mlctrez/godom/example/navbar"
	"strings"
)

var _ app.Handler = (*router)(nil)

type router struct{}

func New() app.Handler {
	return &router{}
}

func (e *router) Prepare(ctx *app.ServerContext) {
	ctx.Main = "example/bin/main.go"
	ctx.Output = "build/app.wasm"
	ctx.Watch = []string{"example", "app"}
	ctx.Address = ":8080"
}

//go:embed head.html
var headHtml string

func (e *router) Headers(ctx *app.Context, header godom.Element) {
	if header.Parent() != nil {
		header.Parent().SetAttribute("data-bs-theme", "dark")
	}
	for i, node := range header.ChildNodes() {
		if node.NodeName() == "title" {
			// TODO: simplify this
			oldTitle := header.ChildNodes()[i]
			header.ChildNodes()[i] = ctx.Doc.H("<title>godom</title>")
			oldTitle.(godom.Element).Remove()
		}
	}
	if len(header.GetElementsByTagName("link")) == 0 {
		header.Body(ctx.Doc.H(headHtml).ChildNodes()...)
	}
}

func (e *router) Body(ctx *app.Context) godom.Element {
	d := ctx.Doc

	body := d.El("body").Body(navbar.Render(ctx))
	switch ctx.URL.Path {
	case "/":
		return body
	case "/exampleOne":
		return body.Body(demos.ExampleOne(ctx))
	case "/diff":
		return body.Body(demos.Diff(ctx))
	case "/editor":
		return body.Body(demos.Editor(ctx))
	default:
		return body.Body(d.H(strings.Replace(was404, "@@page@@", ctx.URL.String(), 1)))
	}
}

//go:embed 404.html
var was404 string
