package main

import (
	"github.com/mlctrez/godom/example"
	"github.com/mlctrez/godom/gsrv"
)

func main() {
	gsrv.Run(&example.Example{})
}
