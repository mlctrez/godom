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

func (n *navbar) form(ctx *app.Context) func(event godom.Value) {
	return func(event godom.Value) {
		event.Call("preventDefault")
		for _, el := range n.Form.GetElementsByTagName("input") {
			// process form data, do a fetch request, etc.
			fmt.Println(el.NodeName(), el.This().Get("value").String())
		}
	}
}

func Render(ctx *app.Context) godom.Element {
	// TODO: simplify making a new doc context with callback
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

	var navItemFmt = `<li class="nav-item"><a role="button" class="nav-link">%s</a></li>`

	// TODO: this is messy and a bit more code than expected. find a way to simplify.
	for _, outer := range navItems {
		// loop vars are bad, escape them
		item := outer

		ni := doc.H(fmt.Sprintf(navItemFmt, item.Name))
		// TODO: add navigation to new location with event listener func
		ni.ChildNodes()[0].AddEventListener("click", func(event godom.Value) {
			u, _ := url.Parse(item.Path)
			ctx.Events <- &app.Location{URL: u}
		})
		if ctx.URL.Path == item.Path {
			// TODO: get element from child nodes without cast
			element := ni.ChildNodes()[0].(godom.Element)
			element.SetAttribute("class", "nav-link active")
			element.SetAttribute("aria-current", "page")
		}
		nav.NavItems.AppendChild(ni)
	}

	nav.Form.AddEventListener("submit", nav.form(ctx))
	return fragment
}
