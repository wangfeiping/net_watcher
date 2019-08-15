package main

import (
	"github.com/wangfeiping/net_watcher/commands"
)

func main() {
	root := commands.NewRootCommand(versioner)
	root.AddCommand(
		// commands.NewStartCommand(nil, true),
		// commands.NewAddCommand(nil, false),
		// commands.NewCallCommand(nil, false),
		commands.NewVersionCommand(versioner))

	if err := root.Execute(); err != nil {
	}
}
