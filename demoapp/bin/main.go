package main

import (
	"github.com/mlctrez/godom/demoapp"
	"github.com/mlctrez/godom/nap"
	"github.com/mlctrez/godom/nap/runner"
)

func main() {
	config := nap.New().WithApp(&demoapp.App{})
	config.PagesPath = "/godom/build"
	config.BuildOutput = "build"
	_ = runner.Run(config)
}
