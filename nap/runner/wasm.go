//go:build wasm

package runner

import (
	"github.com/mlctrez/godom/nap"
	"github.com/mlctrez/godom/nap/runner/wasm"
)

func setup(o *nap.Config) error {
	return o.Setup(false)
}

func run(o *nap.Config) error {
	return wasm.Run(o)
}
