package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/nateph/rcse/pkg/cliconfig"
	"github.com/nateph/rcse/pkg/command"
	"github.com/nateph/rcse/pkg/concurrent"
	"github.com/nateph/rcse/pkg/files"

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
	if err := cliconfig.CheckBaseOptions(baseSettings); err != nil {
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

	shellOptions := cliconfig.JobOptions{
		BaseSettings: baseSettings,
		CommandToRun: commandToRun,
	}

	return executeShell(shellOptions)
}

func executeShell(shellOptions cliconfig.JobOptions) error {
	parsedInventoryFile, err := files.LoadInventory(shellOptions.InventoryFile)
	if err != nil {
		return err
	}

	hosts := parsedInventoryFile.Hosts

	if shellOptions.ListHosts {
		for _, host := range hosts {
			fmt.Println(host)
		}
		return nil
	}

	err = concurrent.Execute(hosts, shellOptions)
	if err != nil {
		return err
	}

	return nil
}
