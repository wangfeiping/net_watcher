package commands

import (
	"github.com/spf13/cobra"
)

// NewAddCommand returns add command
func NewAddCommand(run Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdAdd,
		Short: "Add a service url",
		RunE: func(cmd *cobra.Command, args []string) error {
			run()
			return nil
		},
	}

	cmd.Flags().StringP(FlagURL, "u", "", "service url")
	return cmd
}
