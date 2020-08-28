package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wangfeiping/log"
	"github.com/wangfeiping/net_watcher/commands"
)

func main() {
	defer log.Flush()

	cobra.EnableCommandSorting = false

	rootCmd := commands.NewRootCommand(versioner)
	rootCmd.PersistentFlags().String(log.FlagLogFile, "./logs/watcher.log", "log file path")
	viper.BindPFlag(log.FlagLogFile, rootCmd.PersistentFlags().Lookup(log.FlagLogFile))
	rootCmd.PersistentFlags().Int(log.FlagSize, 10, "log size(MB)")
	viper.BindPFlag(log.FlagSize, rootCmd.PersistentFlags().Lookup(log.FlagSize))

	rootCmd.AddCommand(
		commands.NewStartCommand(starter, true),
		commands.NewAddCommand(addHandler),
		commands.NewCallCommand(callHandler),
		commands.NewVersionCommand(versioner))

	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
