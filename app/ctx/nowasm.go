//go:build !wasm

package ctx

import (
	"github.com/mlctrez/godom/app"
	"github.com/mlctrez/godom/app/nowsm"
)

func Run(h app.Handler) {
	nowsm.Run(h)
}
