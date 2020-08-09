package cmd

import (
	"errors"
	"io"

	"github.com/nateph/rcse/pkg/cliconfig"
	"github.com/nateph/rcse/pkg/concurrent"
	"github.com/spf13/cobra"
)

var (
	sequenceExample = `
	# Run a sequence
	rcse sequence -i ~/inv.yaml -f sequence.yaml

	# Run a sequence as a different user
	rcse sequence -i ~/inv.yaml -f sequence.yaml -u root -p

	# Run a sequence with forks
	rcse sequence -i ~/inv.yaml -f sequence.yaml --forks=10 --failure-limit=2
	`
)

// SequenceOptions is the commandline options for 'sequence' sub command
type SequenceOptions struct {
	BaseOpts     *cliconfig.Options
	SequenceFile string
}

// NewSequenceCommand validates and runs the 'shell' sub command
func NewSequenceCommand(out io.Writer) *cobra.Command {
	o := &SequenceOptions{BaseOpts: &cliconfig.Options{}}
	cmd := &cobra.Command{
		Use:     "sequence",
		Short:   "Execute a shell command",
		Example: sequenceExample,
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

	cmd.Flags().StringVarP(&o.SequenceFile, "sequencefile", "f", "", "the sequence file, in yaml format")
	o.BaseOpts.AddBaseFlags(cmd.Flags())

	return cmd
}

// Validate makes sure provided values and valid Job options
func (s *SequenceOptions) Validate() error {
	if s.SequenceFile == "" {
		return errors.New("no sequence file was found passed. exiting")
	}
	if err := s.BaseOpts.CheckBaseOptions(); err != nil {
		return err
	}
	return nil
}

// Run performs the execution of the 'shell' sub command
func (s *SequenceOptions) Run() error {
	inventory, err := cliconfig.LoadInventory(s.BaseOpts.InventoryFilePath)
	if err != nil {
		return err
	}
	config, err := cliconfig.LoadConfig(s.SequenceFile)
	if err != nil {
		return err
	}
	config.Options = *s.BaseOpts

	err = concurrent.Execute(config, inventory.Hosts...)
	if err != nil {
		return err
	}

	return nil
}
