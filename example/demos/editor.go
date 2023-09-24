package demos

import (
	"bytes"
	_ "embed"
	"encoding/xml"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app"
	"github.com/mlctrez/godom/callback"
)

//go:embed editor.html
var editorHtml string

type editor struct {
	Textarea godom.Element `go:"textarea"`
	Target   godom.Element `go:"target"`
}

func (e *editor) bind(ctx *app.Context) {
	e.Textarea.AddEventListener("change", func(event godom.Value) {
		e.update(ctx)
	})
}

func (e *editor) update(ctx *app.Context) {
	for _, node := range e.Target.ChildNodes() {
		e.Target.RemoveChild(node.This())
	}
	s := e.Textarea.This().Get("value").String()
	dec, err := ctx.Doc.Decode(xml.NewDecoder(bytes.NewBufferString(s)))
	e.Target.AppendChild(ctx.Doc.H(s))
	if err == nil {
		// reformat good data back into the text area
		encoder := godom.NewEncoder(&bytes.Buffer{})
		encoder.Indent("  ")
		s = dec.Marshal(encoder).Xml()
		e.Textarea.This().Set("value", s)
	}
}

func Editor(ctx *app.Context) godom.Element {
	ed := &editor{}
	el := ctx.Doc.WithCallback(callback.Reflect(ed)).H(editorHtml)
	ed.bind(ctx)
	ed.update(ctx)
	return el
}
