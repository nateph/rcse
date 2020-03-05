package cmd

import (
	"errors"
	"io"
	"rcse/cmd/cliconfig"

	"github.com/spf13/cobra"
)

var (
	cliSettings *cliconfig.CliSettings
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
	cliSettings = new(cliconfig.CliSettings)

	cmd.AddCommand(
		newShellCommand(out),
		newVersionCommand(out),
	)

	return cmd
}
