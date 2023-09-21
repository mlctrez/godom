//go:build wasm

package gsrv

import (
	"github.com/mlctrez/godom/gsrv/api"
	"github.com/mlctrez/godom/gsrv/wsm"
)

func Run(h api.Handler) {
	wsm.Run(h)
}
