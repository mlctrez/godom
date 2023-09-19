package main

import (
	_ "embed"
	dom "github.com/mlctrez/godom"
	"github.com/mlctrez/godom/app/base"
	"strings"
)

type App struct {
	base.App
}

//go:embed body.html
var bodyString string

type AppEvent struct {
	el     dom.Element
	event  string
	dataGo string
}

func (e *AppEvent) handleEvent(event dom.Value) {
	dom.Console().Log("%s %o %o", e.dataGo, event, e.el.This())
}

// Run is the main entry point
func Run() {
	c := dom.Console()

	document := dom.Document()
	doc := dom.Document().DocApi()
	doc.CallBack = func(e dom.Element, dataGo []string) {
		for _, val := range dataGo {
			split := strings.Split(val, ".")
			if len(split) != 2 {
				c.Log("invalid data-go attribute %s on %o", val, e.This())
				continue
			}
			ae := &AppEvent{el: e, event: split[1], dataGo: split[0]}
			e.AddEventListener(split[1], ae.handleEvent)
		}
	}

	document.Body().ReplaceWith(doc.H(bodyString))

}

func main() {
	(&App{}).RunMain(Run)
}
