package commands

import (
	"github.com/spf13/cobra"
)

// NewCallCommand returns call command
func NewCallCommand(run Runner) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CmdCall,
		Short: "Request the specified network service",
		RunE: func(cmd *cobra.Command, args []string) error {
			run()
			return nil
		},
	}

	cmd.Flags().StringP(FlagURL, "u", "", "service url")
	return cmd
}
