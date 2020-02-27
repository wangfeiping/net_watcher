package commands

import (
	"github.com/spf13/cobra"
)

// NewStartCommand 创建 start/服务启动 命令
func NewStartCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdStart,
		Short: "Start watch all services",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(run, isKeepRunning)
		},
	}

	cmd.Flags().Int64P(FlegDuration, "d", 30, "The cycle time of the watch task")
	cmd.Flags().StringP(FlagListen, "l", ":9900", "The listening address(ip:port) of exporter")
	return cmd
}
