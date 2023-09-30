package boot

import (
	_ "embed"
	"github.com/mlctrez/godom/nap"
)

//go:embed bootstrap.min.css
var css nap.ResourceBytes

//go:embed bootstrap.bundle.min.js
var js nap.ResourceBytes

func Routes(router nap.Router) {
	router.Resource(nap.Re("/boot/bootstrap.min.css"), "text/css", css.Write)
	router.Resource(nap.Re("/boot/bootstrap.bundle.min.js"), "application/javascript", js.Write)
}
