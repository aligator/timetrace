package main

import (
	"github.com/dominikbraun/timetrace/clockify/clockify"
	"os"

	"github.com/dominikbraun/timetrace/cli"
	"github.com/dominikbraun/timetrace/config"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"
)

var version = "UNDEFINED"

func main() {
	cfg, err := config.FromFile()
	if err != nil {
		out.Warn("%s", err.Error())
	}

	clockifyCfg, err := clockify.ConfigFromFile()
	if err != nil {
		out.Warn("%s", err.Error())
	}

	filesystem := clockify.NewFs(clockifyCfg)
	timetrace := core.New(cfg, filesystem)

	if err := cli.RootCommand(timetrace, version).Execute(); err != nil {
		out.Err("%s", err.Error())
		os.Exit(1)
	}
}
