package commands

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wangfeiping/log"
)

// nolint
const (
	CmdRoot          = "net_watcher"
	CmdStart         = "start"
	CmdAdd           = "add"
	CmdCall          = "call"
	CmdConfig        = "config"
	CmdVersion       = "version"
	CmdHelp          = "help"
	ShortDescription = "Network service availability detection tools"
)

// nolint
const (
	FlagLog      = "log"
	FlagConfig   = "config"
	FlagListen   = "listen"
	FlagAlias    = "alias"
	FlagURL      = "url"
	FlagBody     = "body"
	FlagMethod   = "method"
	FlagRegex    = "regex"
	FlegDuration = "duration"
	FlagService  = "service"
	FlagVersion  = CmdVersion
)

// Runner is command call function
type Runner func() (context.CancelFunc, error)

// NewRootCommand returns root command
func NewRootCommand(versioner Runner) *cobra.Command {
	root := &cobra.Command{
		Use:   CmdRoot,
		Short: ShortDescription,
		Run: func(cmd *cobra.Command, args []string) {
			if viper.GetBool(FlagVersion) {
				versioner()
				return
			}
			cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if err := viper.BindPFlags(cmd.Flags()); err != nil {
				log.Error("bind flags error: ", err)
				return err
			}

			if strings.EqualFold(cmd.Use, CmdRoot) ||
				strings.EqualFold(cmd.Use, CmdVersion) {
				// doesn't need init config & log
				return nil
			}

			initConfig()

			if !strings.EqualFold(cmd.Use, CmdStart) {
				// doesn't need init log
				return nil
			}

			initLogger()

			return
		},
	}

	root.Flags().BoolP(FlagVersion, "v", false, "Show version info")
	root.PersistentFlags().StringP(FlagConfig, "c", "./config.yml", "Config file path")

	return root
}

func initConfig() error {
	viper.SetConfigFile(viper.GetString(FlagConfig))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// stderr, so if we redirect output to json file, this doesn't appear
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		// ignore not found error, return other errors
		return err
	}

	return nil
}

func initLogger() {
	log.Config(log.RollingFileConfig())
	log.Infof("starting at %s", getExecPath())
}

// getExecPath returns the execution path
func getExecPath() (execPath string) {
	file, _ := exec.LookPath(os.Args[0])
	execFile := filepath.Base(file)
	execPath, _ = filepath.Abs(file)
	if len(execPath) > 1 {
		rs := []rune(execPath)
		execPath = string(rs[0:(len(execPath) - len(execFile))])
	}
	return
}

func commandRunner(run Runner, isKeepRunning bool) error {
	cancel, err := run()
	if err != nil {
		log.Error("Running error: ", err.Error())
		return err
	}
	if isKeepRunning {
		keepRunning(func(sig os.Signal) {
			defer log.Flush()
			if cancel != nil {
				cancel()
			}
			log.Debug("Stopped by signal: ", sig)
		})
	}
	return nil
}

func keepRunning(callback func(sig os.Signal)) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case s, ok := <-signals:
		log.Infof("System signal [%v] %t, trying to run callback...", s, ok)
		if !ok {
			break
		}
		if callback != nil {
			callback(s)
		}
		log.Flush()
		os.Exit(1)
	}
}

func serviceFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(FlagURL, "u", "", "request url")
	cmd.Flags().StringP(FlagBody, "b", "", "request body")
	cmd.Flags().StringP(FlagAlias, "a", "", "service alias")
	cmd.Flags().StringP(FlagMethod, "m", "GET", "http method")
	cmd.Flags().StringP(FlagRegex, "r", "GET", "match regex")

}
