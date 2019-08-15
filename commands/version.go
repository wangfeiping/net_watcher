package commands

import (
	"github.com/spf13/cobra"
)

// NewVersionCommand returns version command
func NewVersionCommand(run Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdVersion,
		Short: "Show version info",
		RunE: func(cmd *cobra.Command, args []string) error {
			run()
			return nil
		},
	}
	return cmd
}
