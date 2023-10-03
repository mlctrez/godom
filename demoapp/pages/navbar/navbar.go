package navbar

import (
	_ "embed"
	"fmt"
	"github.com/mlctrez/godom/nap"
	"net/url"
	"strings"

	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/callback"
)

//go:embed navbar.html
var navbarHtml string

type Navbar struct {
	// NavItems is line 10 in navbar.html
	NavItems godom.Element `go:"navigationItems"`
	// Form is line 12 in navbar.html
	Form godom.Element `go:"searchForm"`
}

func (n *Navbar) searchSubmit(event godom.Value) {
	event.Call("preventDefault")
	for _, el := range n.Form.GetElementsByTagName("input") {
		// this would really do a form post
		searchString := el.This().Get("value").String()
		if strings.TrimSpace(searchString) != "" {
			fmt.Printf("searching for %q\n", searchString)
		}
	}
}

func Render(ctx nap.DocContext) godom.Element {
	nav := &Navbar{}

	api := ctx.DocApi()
	fragment := api.WithCallback(callback.Reflect(nav)).H(navbarHtml)
	nav.finalize(ctx)
	return fragment
}

func (n *Navbar) finalize(ctx nap.DocContext) {
	if n.Form == nil || n.NavItems == nil {
		// mismatch between *Navbar struct and html attributes, bail now
		// or nil pointer errors at code locations would probably be enough to figure out what was broken
		panic("bind Form and NavItems failed")
	}
	n.Form.AddEventListener("submit", n.searchSubmit)
	n.buildNavItems(ctx)
}

func (n *Navbar) buildNavItems(ctx nap.DocContext) {
	// TODO: this is messy and a bit more code than expected. find a way to simplify.
	for _, outer := range navItems {
		// escape loop var
		item := outer

		cb := func(e godom.Element, name, data string) {
			if (item.Path == "/" && ctx.URL().Path == "/index") || ctx.URL().Path == item.Path {
				e.SetAttribute("class", "nav-link active")
				e.SetAttribute("aria-current", "page")
			} else {
				e.AddEventListener("click", func(event godom.Value) {
					u, _ := url.Parse(item.Path)
					ctx.Event(&nap.Location{URL: u})
				})
			}
		}

		var navItemFmt = `<li class="nav-item"><a go="anchor" role="button" class="nav-link">%s</a></li>`
		ni := ctx.DocApi().WithCallback(cb).H(fmt.Sprintf(navItemFmt, item.Name))
		n.NavItems.AppendChild(ni)
	}
}

var navItems = []struct {
	Name string
	Path string
}{
	{Name: "Home", Path: "/"},
	{Name: "About", Path: "/about"},
}
