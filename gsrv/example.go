package gsrv

import (
	"github.com/mlctrez/godom"
	"github.com/mlctrez/godom/gsrv/api"
	"net/url"
)

var _ api.Handler = (*Example)(nil)

type Example struct {
}

func (e *Example) Body(doc godom.Doc, u *url.URL) godom.Element {
	return doc.H("<body><p>" + u.String() + "</p></body>")
}
