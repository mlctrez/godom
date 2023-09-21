//go:build !wasm

package gsrv

import (
	"github.com/mlctrez/godom/gsrv/api"
	"github.com/mlctrez/godom/gsrv/nowsm"
)

func Run(h api.Handler) {
	nowsm.Run(h)
}
