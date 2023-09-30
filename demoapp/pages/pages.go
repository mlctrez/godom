package pages

import (
	_ "embed"
	"fmt"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/demoapp/pages/boot"
	"github.com/mlctrez/godom/nap"
	"github.com/mlctrez/wasmexec"
	"io"
	"net/url"
)

type Pages struct {
	c *nap.Config
}

func New(c *nap.Config) *Pages {
	return &Pages{c: c}
}

func (p *Pages) Routes(router nap.Router) {
	router.Page(nap.Re("/$|^/index"), p.pageFunc(&indexHtml))
	router.Page(nap.Re("/about"), p.pageFunc(&aboutHtml))
	router.Resource(nap.Re("/app.js"), "application/javascript", p.appJS)
	router.Resource(nap.Re("/app.wasm"), "application/wasm", wasm.Write)
	boot.Routes(router)
	router.NotFound(func(u *url.URL, ctx nap.DocContext) godom.Element {
		body := fmt.Sprintf("<body>404: page %s not found</body>", u.Path)
		return p.pageFunc(&body)(ctx)
	})
}

//go:embed app.js
var appjs []byte

func (p *Pages) appJS(writer io.Writer) (n int, err error) {
	var content []byte
	if content, err = wasmexec.Current(); err != nil {
		return 0, err
	}
	content = append(content, appjs...)
	return writer.Write(content)
}

//go:embed app.wasm
var wasm nap.ResourceBytes

//go:embed head.html
var headHtml string

//go:embed index.html
var indexHtml string

func (p *Pages) pageFunc(body *string) nap.PageFunc {
	return func(ctx nap.DocContext) (page godom.Element) {
		if p.c.IsServer() {
			page = ctx.DocApi().H(headHtml)
		}
		bodyEl := ctx.DocApi().H(*body)
		if page != nil {
			page.AppendChild(bodyEl)
		} else {
			page = bodyEl
		}
		return page
	}
}

//func (p *Pages) index(ctx nap.DocContext) godom.Element {
//	html := ctx.DocApi().H(headHtml)
//	html.AppendChild(ctx.DocApi().H(indexHtml))
//	return html
//}

//go:embed about.html
var aboutHtml string

//func (p *Pages) about(ctx nap.DocContext) godom.Element {
//	return ctx.DocApi().H(aboutHtml)
//}
