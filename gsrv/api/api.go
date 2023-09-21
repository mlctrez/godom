package api

import (
	"github.com/mlctrez/godom"
	"net/url"
)

type Handler interface {
	Body(doc godom.Doc, u *url.URL) godom.Element
}
