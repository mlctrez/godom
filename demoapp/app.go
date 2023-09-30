package demoapp

import (
	_ "embed"
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/demoapp/pages"
	"github.com/mlctrez/godom/nap"
)

var _ nap.App = (*App)(nil)

type App struct {
	c *nap.Config
	h godom.Element
}

func (a *App) Init(c *nap.Config, html godom.Element) {
	a.c = c
	a.h = html
}

func (a *App) Routes(router nap.Router) {
	pages.New(a.c).Routes(router)
}
