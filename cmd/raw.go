package cmd

import (
	"fmt"
	"rcse/cmd/cliconfig"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	commandToRun string
	rawShellCmd  = &cobra.Command{
		Use:   "raw",
		Short: "Execute a raw shell command",
		Long:  "Execute a raw shell command on a remote host",
		Run:   rawCommand,
	}
)

func init() {
	rootCmd.AddCommand(rawShellCmd)
	rawShellCmd.Flags().StringVarP(&commandToRun, "command", "c", "", "the command to run on a remote host")
	rawShellCmd.MarkFlagRequired("command")
}

func rawCommand(cmd *cobra.Command, args []string) {
	parsedHosts := viper.GetStringSlice("hosts")
	// listHosts is a persistent flag set in the root command
	if listHosts {
		for _, host := range parsedHosts {
			fmt.Println(host)
		}
		return
	}
	var wg sync.WaitGroup

	for _, host := range parsedHosts {
		wg.Add(1)
		rawCmdOpts := cliconfig.CommandOptions{
			Host:               host,
			CommandToRun:       commandToRun,
			Sudo:               false,
			IgnoreHostkeyCheck: ignoreHostkeyCheck,
		}
		go func(host string) {
			defer wg.Done()
			rawCmdOpts.RunCommand()
		}(host)
	}
	wg.Wait()
}
