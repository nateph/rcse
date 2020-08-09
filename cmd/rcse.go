package cmd

import (
	"errors"
	"io"

	"github.com/spf13/cobra"
)

// NewRcseCommand returns a new rcse command
func NewRcseCommand(out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "rcse",
		Short:         "Run a command somewhere else",
		Long:          "Run a command somewhere else",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return errors.New("no arguments accepted")
			}
			return nil
		},
	}

	cmd.PersistentFlags().Parse(args)

	cmd.AddCommand(
		NewSequenceCommand(out),
		NewShellCommand(out),
		NewVersionCommand(out),
	)

	return cmd
}
