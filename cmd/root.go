package cmd

import (
	"errors"
	"io"

	"github.com/nateph/rcse/pkg/cliconfig"

	"github.com/spf13/cobra"
)

var (
	baseSettings *cliconfig.Options
)

// NewRootCmd returns a root command
func NewRootCmd(out io.Writer, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rcse",
		Short: "Run a command somewhere else",
		Long:  "Run a command somewhere else",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return errors.New("no arguments accepted")
			}
			return nil
		},
	}

	flags := cmd.PersistentFlags()
	flags.Parse(args)
	baseSettings = new(cliconfig.Options)

	cmd.AddCommand(
		newSequenceCommand(out),
		newShellCommand(out),
		newVersionCommand(out),
	)

	return cmd
}
