package cmd

import (
	"fmt"
	"io"
	"rcse/cmd/cliconfig"
	"rcse/pkg/files"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	commandToRun string
)

// RawOptions contains options for the raw command
type RawOptions struct {
	CommandToRun       string
	IgnoreHostKeyCheck bool
	InventoryFile      string
	ListHosts          bool
	Password           string
	User               string
}

func newRawCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "raw",
		Short:        "Execute a raw shell command",
		SilenceUsage: true,
		RunE:         runRaw,
	}

	flags := cmd.Flags()
	cliSettings.AddFlags(flags)
	flags.StringVarP(&commandToRun, "command", "c", "", "the command to run on a remote host")

	return cmd
}

func runRaw(cmd *cobra.Command, args []string) error {
	if cliSettings.InventoryFile == "" {
		logrus.Fatal("no inventory flag was specified, all rcse operations require an inventory")
	}
	if commandToRun == "" {
		logrus.Fatal("No command was found to run. Exiting.")
	}
	// if --username and --password were supplied correctly without --list-hosts
	if cliSettings.User != "" && cliSettings.Password == "default" && !cliSettings.ListHosts {
		cliSettings.Password = cliconfig.CheckAndConsumePassword(cliSettings.User, cliSettings.Password)
	}

	rawOptions := RawOptions{
		IgnoreHostKeyCheck: cliSettings.IgnoreHostKeyCheck,
		InventoryFile:      cliSettings.InventoryFile,
		ListHosts:          cliSettings.ListHosts,
		Password:           cliSettings.Password,
		User:               cliSettings.User,
		CommandToRun:       commandToRun,
	}

	return executeRaw(rawOptions)
}

func executeRaw(rawOptions RawOptions) error {
	parsedInventoryFile, err := files.LoadInventory(rawOptions.InventoryFile)
	if err != nil {
		return err
	}

	parsedHosts := parsedInventoryFile.Hosts

	if rawOptions.ListHosts {
		for _, host := range parsedHosts {
			fmt.Println(host)
		}
		return nil
	}
	// ---------------------------------------------
	results := make(chan cliconfig.CommandResult)
	timeout := time.After(10 * time.Second)

	for _, host := range parsedHosts {
		rawCmdOpts := cliconfig.CommandOptions{
			Host:               host,
			CommandToRun:       commandToRun,
			Sudo:               false,
			IgnoreHostkeyCheck: rawOptions.IgnoreHostKeyCheck,
			User:               rawOptions.User,
			Password:           rawOptions.Password,
		}
		go func() {
			results <- rawCmdOpts.RunCommands()
		}()
	}

	for i := 0; i < len(parsedHosts); i++ {
		select {
		case res := <-results:
			res.PrintHostOutput()
		case <-timeout:
			fmt.Println("timed out")
		}
	}
	return nil
}
