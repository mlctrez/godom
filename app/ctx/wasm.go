//go:build wasm

package ctx

import (
	"github.com/mlctrez/godom/app"
	"github.com/mlctrez/godom/app/wsm"
)

func Run(h app.Handler) {
	wsm.Run(h)
}
