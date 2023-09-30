//go:build !wasm

package runner

import (
	"github.com/mlctrez/godom/nap"
	"github.com/mlctrez/godom/nap/runner/server"
)

func setup(o *nap.Config) error {
	return o.Setup(true)
}

func run(o *nap.Config) error {
	return server.Run(o)
}
