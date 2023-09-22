package example

import (
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app"
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
		return e.notFound(ctx)
	}
}

func (e *example) index(ctx *app.Context) godom.Element {
	body := ctx.Doc.H("<body><table><tbody/></table></body>")
	body.GetElementsByTagName("tbody")[0].AppendChild((&exampleOne{}).render(ctx))
	return body
}

func (e *example) notFound(ctx *app.Context) godom.Element {
	return ctx.Doc.H(pageNotFound)
}

var pageNotFound = `<body>
<p style="color:red">page not found</p>
<br/>
Go back to <a href="/">index page</a>
</body>`
