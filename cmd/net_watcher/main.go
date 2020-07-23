package main

import (
	"github.com/wangfeiping/log"
	"github.com/wangfeiping/net_watcher/commands"
)

func main() {
	defer log.Flush()

	root := commands.NewRootCommand(versioner)
	root.AddCommand(
		commands.NewStartCommand(starter, true),
		commands.NewAddCommand(addHandler),
		commands.NewCallCommand(callHandler),
		commands.NewVersionCommand(versioner))

	if err := root.Execute(); err != nil {
	}
}
