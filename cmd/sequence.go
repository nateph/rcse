package cmd

import (
	"fmt"
	"io"

	"github.com/nateph/rcse/pkg/cliconfig"
	"github.com/nateph/rcse/pkg/concurrent"
	"github.com/spf13/cobra"
)

var (
	sequenceFile string
)

func newSequenceCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "sequence",
		Short:        "Run a sequence of jobs",
		SilenceUsage: true,
		RunE:         runSequence,
	}

	flags := cmd.Flags()
	baseSettings.AddFlags(flags)
	flags.StringVarP(&sequenceFile, "file", "f", "", "the sequence file, in yaml format")

	return cmd
}

func runSequence(cmd *cobra.Command, args []string) error {
	if err := baseSettings.CheckBaseOptions(); err != nil {
		return err
	}

	inventory, err := cliconfig.LoadInventory(baseSettings.InventoryFilePath)
	if err != nil {
		return err
	}

	config, err := cliconfig.LoadConfig(sequenceFile)
	if err != nil {
		return err
	}
	config.Options = *baseSettings

	return executeSequence(config, inventory.Hosts...)
}

func executeSequence(config *cliconfig.Config, inventory ...string) error {
	if config.Options.ListHosts {
		for _, host := range inventory {
			fmt.Println(host)
		}
		return nil
	}

	err := concurrent.Execute(config, inventory...)
	if err != nil {
		return err
	}

	return nil
}
