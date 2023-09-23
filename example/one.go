package example

import (
	"fmt"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app"
	"github.com/mlctrez/godom/callback"
	"os"
	"time"
)

type exampleOne struct {
	Button godom.Element `go:"button"`
	Reset  godom.Element `go:"reset"`
	Div    godom.Element `go:"div"`
}

var exOneRow = `<tr><td>
<button go="button" type="button">example one</button><button go="reset">reset</button> 
<div go="div"></div>
</td></tr>`

func (eo *exampleOne) callBack() func(e godom.Element, name string, data string) {
	// callBack() can use one or the other
	if os.Getenv("NO_REFLECTION") != "" {
		// with tags and reflect
		return callback.Reflect(eo)
	} else {
		// without tags or reflection
		return callback.Mapper(map[string]func(godom.Element){
			"button": func(ei godom.Element) { eo.Button = ei },
			"reset":  func(ei godom.Element) { eo.Reset = ei },
			"div":    func(ei godom.Element) { eo.Div = ei },
		})
	}
}

func (eo *exampleOne) events(doc godom.Doc) {
	eo.Button.AddEventListener("click", func(event godom.Value) {
		//eo.Button.SetAttribute("disabled", true)
		//eo.Button.This().Set("innerHTML", "disabled")
		//go time.AfterFunc(1*time.Second, func() {
		//	eo.Button.RemoveAttribute("disabled")
		//	eo.Button.This().Set("innerHTML", "example one")
		//})
		eo.Div.AppendChild(doc.El("br"))
		eo.Div.AppendChild(doc.H(fmt.Sprintf("<span>%s</span>", time.Now().Format(time.RFC3339Nano))))
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
}

func (eo *exampleOne) render(ctx *app.Context) godom.Element {
	doc := godom.Doc{Doc: ctx.Doc.Doc, CallBack: eo.callBack()}
	row := doc.H(exOneRow)
	eo.events(doc)
	return row
}
