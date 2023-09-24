package demos

import (
	_ "embed"
	"fmt"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app"
	"github.com/mlctrez/godom/callback"
	"time"
)

type exampleOne struct {
	Button godom.Element `go:"button"`
	Reset  godom.Element `go:"reset"`
	Div    godom.Element `go:"div"`
}

func ExampleOne(ctx *app.Context) godom.Element {
	eo := &exampleOne{}
	doc := ctx.Doc.WithCallback(callback.Reflect(eo))
	row := doc.H(exOneHtml)
	eo.Button.AddEventListener("click", func(event godom.Value) {

		span := fmt.Sprintf("<span>%s</span>", time.Now().Format(time.RFC3339Nano))
		eo.Div.Body(doc.El("br"), doc.H(span))

		if len(eo.Div.ChildNodes()) > 12 {
			eo.Div.RemoveChild(eo.Div.ChildNodes()[0].This())
			eo.Div.RemoveChild(eo.Div.ChildNodes()[0].This())
		}
	})
	eo.Reset.AddEventListener("click", func(event godom.Value) {
		for len(eo.Div.ChildNodes()) > 0 {
			eo.Div.RemoveChild(eo.Div.ChildNodes()[0].This())
		}
	})
	return row
}

//go:embed one.html
var exOneHtml string
