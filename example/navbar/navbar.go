package navbar

import (
	_ "embed"
	"fmt"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app"
	"github.com/mlctrez/godom/callback"
	"net/url"
)

//go:embed navbar.html
var navbarHtml string

type navbar struct {
	NavItems godom.Element `go:"navigationItems"`
	Form     godom.Element `go:"searchForm"`
}

func (n *navbar) nav(ctx *app.Context, u string) func(event godom.Value) {
	return func(event godom.Value) {
		url, _ := url.Parse(u)
		ctx.Events <- &app.Location{URL: url}
	}
}

func (n *navbar) form(ctx *app.Context) func(event godom.Value) {
	return func(event godom.Value) {
		event.Call("preventDefault")
		for _, el := range n.Form.GetElementsByTagName("input") {
			// process form data, do a fetch request, etc.
			fmt.Println(el.NodeName(), el.This().Get("value").String())
		}
	}
}

type navItem struct {
	Element godom.Element `go:"element"`
}

func Render(ctx *app.Context) godom.Element {
	doc := godom.Doc{Doc: godom.Document().This()}
	nav := &navbar{}
	doc.CallBack = callback.Reflect(nav)
	fragment := doc.H(navbarHtml)

	navItems := []struct {
		Name string
		Path string
	}{
		{Name: "Home", Path: "/"},
		{Name: "Page", Path: "/page"},
	}

	for _, outer := range navItems {
		item := outer
		li := doc.El("li", doc.At("class", "nav-item"))
		a := doc.El("a", doc.At("role", "button"))
		li.AppendChild(a)
		a.AddEventListener("click", func(event godom.Value) {
			u, _ := url.Parse(item.Path)
			fmt.Println("click", u)
			ctx.Events <- &app.Location{URL: u}
		})
		if ctx.URL.Path == item.Path {
			a.SetAttribute("class", "nav-link active")
			a.SetAttribute("aria-current", "page")
		} else {
			a.SetAttribute("class", "nav-link")
		}
		a.AppendChild(doc.T(item.Name))
		nav.NavItems.AppendChild(li)
	}

	nav.Form.AddEventListener("submit", nav.form(ctx))
	return fragment
}
