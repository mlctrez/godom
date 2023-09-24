package navbar

import (
	_ "embed"
	"fmt"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app"
	"github.com/mlctrez/godom/callback"
	"net/url"
	"strings"
)

//go:embed navbar.html
var navbarHtml string

type navbar struct {
	// NavItems is line 10 in navbar.html
	NavItems godom.Element `go:"navigationItems"`
	// Form is line 12 in navbar.html
	Form godom.Element `go:"searchForm"`
}

func (n *navbar) searchSubmit(event godom.Value) {
	event.Call("preventDefault")
	for _, el := range n.Form.GetElementsByTagName("input") {
		// this would really do a form post
		searchString := el.This().Get("value").String()
		if strings.TrimSpace(searchString) != "" {
			fmt.Printf("searching for %q\n", searchString)
		}
	}
}

func Render(ctx *app.Context) godom.Element {
	nav := &navbar{}
	fragment := ctx.Doc.WithCallback(callback.Reflect(nav)).H(navbarHtml)
	nav.finalize(ctx)
	return fragment
}

func (n *navbar) finalize(ctx *app.Context) {
	if n.Form == nil || n.NavItems == nil {
		// mismatch between *navbar struct and html attributes, bail now
		// or nil pointer errors at code locations would probably be enough to figure out what was broken
		panic("bind Form and NavItems failed")
	}
	n.Form.AddEventListener("submit", n.searchSubmit)
	n.buildNavItems(ctx)
}

func (n *navbar) buildNavItems(ctx *app.Context) {

	// TODO: this is messy and a bit more code than expected. find a way to simplify.
	for _, outer := range navItems {
		// escape loop var
		item := outer

		cb := func(e godom.Element, name, data string) {
			//fmt.Println(ctx.URL.Path, item.Path, e)
			if ctx.URL.Path == item.Path {
				e.SetAttribute("class", "nav-link active")
				e.SetAttribute("aria-current", "page")
			}
			e.AddEventListener("click", func(event godom.Value) {
				u, _ := url.Parse(item.Path)
				ctx.Events <- &app.Location{URL: u}
			})
		}

		var navItemFmt = `<li class="nav-item"><a go="anchor" role="button" class="nav-link">%s</a></li>`
		ni := ctx.Doc.WithCallback(cb).H(fmt.Sprintf(navItemFmt, item.Name))
		n.NavItems.AppendChild(ni)
	}
}

var navItems = []struct {
	Name string
	Path string
}{
	{Name: "Home", Path: "/"},
	{Name: "Example One", Path: "/exampleOne"},
	{Name: "Editor", Path: "/editor"},
	{Name: "Diff", Path: "/diff"},
}
