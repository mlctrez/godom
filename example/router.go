package example

import (
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app"
	"github.com/mlctrez/godom/example/navbar"
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

func (e *router) Headers(ctx *app.Context, header godom.Element) {
	if header.Parent() != nil {
		header.Parent().SetAttribute("data-bs-theme", "dark")
	}
	for i, node := range header.ChildNodes() {
		if node.NodeName() == "title" {
			// TODO: this needs to be a replace with
			oldTitle := header.ChildNodes()[i]
			header.ChildNodes()[i] = ctx.Doc.H("<title>godom</title>")
			oldTitle.(godom.Element).Remove()
		}
	}
	if len(header.GetElementsByTagName("link")) == 0 {
		header.AppendChild(ctx.Doc.H(bootstrapCss))
		header.AppendChild(ctx.Doc.H(bootstrapJs))
	}
}

var bootstrapCss = `<link 
href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" 
rel="stylesheet" integrity="sha384-T3c6CoIi6uLrA9TneNEoa7RxnatzjcDSCmG1MXxSR1GAsXEV/Dwwykc2MPK8M2HN" 
crossorigin="anonymous"/>`

var bootstrapJs = `<script 
src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js" 
integrity="sha384-C6RzsynM9kWDrMNeT87bh95OGNyZPhcTNXj1NW7RuBCsyN/o0jlpcV8Qyq46cDfL" 
crossorigin="anonymous"></script>`

func (e *router) Body(ctx *app.Context) godom.Element {
	switch ctx.URL.Path {
	case "/":
		return e.index(ctx)
	case "/page":
		return e.index(ctx)
	default:
		return e.notFound(ctx)
	}
}

func (e *router) index(ctx *app.Context) godom.Element {

	doc := ctx.Doc
	body := doc.El("body")

	body.AppendChild(navbar.Render(ctx))
	// TODO: remove example one since it is not used
	//body := doc.H("<body><table><tbody/></table></body>")
	//body.GetElementsByTagName("tbody")[0].AppendChild((&exampleOne{}).render(ctx))
	return body
}

func (e *router) notFound(ctx *app.Context) godom.Element {
	return ctx.Doc.H(pageNotFound)
}

var pageNotFound = `<body>
<p style="color:red">page not found</p>
<br/>
Go back to <a href="/">index page</a>
</body>`
