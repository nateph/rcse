package cmd

import (
	"io"

	"github.com/nateph/rcse/pkg/cliconfig"
	"github.com/nateph/rcse/pkg/concurrent"
	"github.com/nateph/rcse/pkg/files"
	"github.com/nateph/rcse/pkg/files/inventory"

	"github.com/spf13/cobra"
)

var (
	scriptExample = `
# Run a script
rcse script my_script.sh -i ~/inv.yaml
rcse script -i ~/inv.yaml my_script.sh

# Run a script as a different user
rcse script my_script.sh -i ~/inv.yaml -u root -p

# Run a script with forks
rcse script my_script.sh -i ~/inv.yaml --forks=10 --failure-limit=2
`
)

// ScriptOptions is the commandline options for 'shell' sub command
type ScriptOptions struct {
	BaseOpts       *cliconfig.Options
	ScriptFilePath string
	Cleanup        bool
}

// NewScriptCommand validates and runs the 'shell' sub command
func NewScriptCommand(out io.Writer) *cobra.Command {
	o := &ScriptOptions{BaseOpts: &cliconfig.Options{}}
	cmd := &cobra.Command{
		Use:     "script",
		Short:   "Execute a script remotely",
		Example: scriptExample,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			o.ScriptFilePath = args[0]
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}
			return nil
		},
	}
	o.BaseOpts.AddBaseFlags(cmd.Flags())
	cmd.Flags().BoolVar(&o.Cleanup, "cleanup", false, "delete the script on the remote host after executing it")
	return cmd
}

// Validate makes sure provided values and valid Job options
func (s *ScriptOptions) Validate() error {
	if err := s.BaseOpts.CheckBaseOptions(); err != nil {
		return err
	}

	var err error
	s.ScriptFilePath, err = files.ParseAndVerifyFilePath(s.ScriptFilePath)
	if err != nil {
		return err
	}

	err = files.VerifyScript(s.ScriptFilePath)
	if err != nil {
		return err
	}

	return nil
}

// Run performs the execution of the 'script' sub command
func (s *ScriptOptions) Run() error {
	inventory, err := inventory.LoadInventory(s.BaseOpts.InventoryFilePath)
	if err != nil {
		return err
	}
	job := &cliconfig.Job{
		Script:  s.ScriptFilePath,
		Cleanup: s.Cleanup,
	}
	executeConfig := &cliconfig.Config{
		Job:     *job,
		Options: *s.BaseOpts,
	}
	err = concurrent.Execute(executeConfig, inventory...)
	if err != nil {
		return err
	}

	return nil
}
