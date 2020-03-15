package cmd

import (
	"errors"
	"io"

	"github.com/nateph/rcse/pkg/cliconfig"
	"github.com/nateph/rcse/pkg/command"
	"github.com/nateph/rcse/pkg/concurrent"

	"github.com/spf13/cobra"
)

var (
	commandToRun string
)

func newShellCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "shell",
		Short:        "Execute a shell command",
		SilenceUsage: true,
		RunE:         runShell,
	}

	flags := cmd.Flags()
	baseSettings.AddFlags(flags)
	flags.StringVarP(&commandToRun, "command", "c", "", "the command to run on a remote host")

	return cmd
}

func runShell(cmd *cobra.Command, args []string) error {
	if err := baseSettings.CheckBaseOptions(); err != nil {
		return err
	}

	if commandToRun == "" {
		return errors.New("no command was found to run. exiting")
	}

	// if --username and --password were supplied correctly without --list-hosts
	var err error
	if baseSettings.User != "" && baseSettings.Password == "default" && !baseSettings.ListHosts {
		baseSettings.Password, err = command.CheckAndConsumePassword(baseSettings.User, baseSettings.Password)
		if err != nil {
			return err
		}
	}

	inventory, err := cliconfig.LoadInventory(baseSettings.InventoryFilePath)
	if err != nil {
		return err
	}

	shellConfig := generateShellConfig(*baseSettings)

	return ExecuteShell(shellConfig, inventory.Hosts...)
}

// ExecuteShell executes shell commands concurrently
func ExecuteShell(config *cliconfig.Config, inventory ...string) error {
	err := concurrent.Execute(config, inventory...)
	if err != nil {
		return err
	}

	return nil
}

func generateShellConfig(opts cliconfig.Options) *cliconfig.Config {
	var jobs []cliconfig.Job
	shellConfigJobs := cliconfig.Job{
		Command: commandToRun,
	}
	jobs = append(jobs, shellConfigJobs)

	shellConfig := &cliconfig.Config{
		Jobs:    jobs,
		Options: opts,
	}

	return shellConfig
}
