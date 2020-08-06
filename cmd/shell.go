package cmd

import (
	"errors"
	"io"

	"github.com/nateph/rcse/pkg/cliconfig"
	"github.com/nateph/rcse/pkg/concurrent"

	"github.com/spf13/cobra"
)

var (
	shellExample = `
	# Run a command 
	rcse shell -i ~/inv.yaml -c "ls -la"

	# Run a command as a different user
	rcse shell -i ~/inv.yaml -c "systemctl restart nginx" -u root -p

	# Run a command with forks
	rcse shell -i ~/inv.yaml -c "ls -la" --forks=10 --failure-limit=2
	`
)

// ShellOptions is the commandline options for 'shell' sub command
type ShellOptions struct {
	BaseOpts *cliconfig.Options
	Command  string
}

// NewShellCommand validates and runs the 'shell' sub command
func NewShellCommand(out io.Writer) *cobra.Command {
	o := &ShellOptions{BaseOpts: &cliconfig.Options{}}
	cmd := &cobra.Command{
		Use:     "shell",
		Short:   "Execute a shell command",
		Example: shellExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&o.Command, "command", "c", "", "the command to run on a remote host")
	o.BaseOpts.AddBaseFlags(cmd.Flags())

	return cmd
}

// Validate makes sure provided values and valid Job options
func (s *ShellOptions) Validate() error {
	if err := s.BaseOpts.CheckBaseOptions(); err != nil {
		return err
	}
	if s.Command == "" {
		return errors.New("no command was found to run. exiting")
	}
	return nil
}

// Run performs the execution of the 'shell' sub command
func (s *ShellOptions) Run() error {
	inventory, err := cliconfig.LoadInventory(s.BaseOpts.InventoryFilePath)
	if err != nil {
		return err
	}

	job := &cliconfig.Job{
		Command: s.Command,
	}
	executeConfig := &cliconfig.Config{
		Jobs:    []cliconfig.Job{*job},
		Options: *s.BaseOpts,
	}
	err = concurrent.Execute(executeConfig, inventory.Hosts...)
	if err != nil {
		return err
	}

	return nil
}
