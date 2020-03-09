package cmd

import (
	"fmt"
	"io"

	"github.com/nateph/rcse/pkg/cliconfig"
	"github.com/nateph/rcse/pkg/concurrent"
	"github.com/nateph/rcse/pkg/files"
	"github.com/spf13/cobra"
)

func newPingCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "ping",
		Short:        "Ping hosts from your local machine to ensure connectivity",
		SilenceUsage: true,
		RunE:         runPing,
	}

	flags := cmd.Flags()
	baseSettings.AddFlags(flags)

	return cmd
}

func runPing(cmd *cobra.Command, args []string) error {
	if err := cliconfig.CheckBaseOptions(baseSettings); err != nil {
		return err
	}

	pingOptions := cliconfig.JobOptions{
		BaseSettings: baseSettings,
	}

	return executePing(pingOptions)
}

func executePing(pingOptions cliconfig.JobOptions) error {
	parsedInventoryFile, err := files.LoadInventory(pingOptions.InventoryFile)
	if err != nil {
		return err
	}

	hosts := parsedInventoryFile.Hosts

	if pingOptions.ListHosts {
		for _, host := range hosts {
			fmt.Println(host)
		}
		return nil
	}

	err = concurrent.Execute(hosts, pingOptions)
	if err != nil {
		return err
	}

	return nil
}
