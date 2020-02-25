package cmd

import (
	"fmt"
	"rcse/cmd/cliconfig"
	"time"

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
	// listHostsFlag is a persistent flag set in the root command
	if listHostsFlag {
		for _, host := range parsedHosts {
			fmt.Println(host)
		}
		return
	}
	results := make(chan cliconfig.CommandResult)
	timeout := time.After(10 * time.Second)

	for _, host := range parsedHosts {
		// ignoreHostkeyCheck is a persistent flag set in the root command
		rawCmdOpts := cliconfig.CommandOptions{
			Host:               host,
			CommandToRun:       commandToRun,
			Sudo:               false,
			IgnoreHostkeyCheck: ignoreHostkeyCheck,
		}
		go func() {
			results <- rawCmdOpts.RunCommand()
		}()
	}

	for i := 0; i < len(parsedHosts); i++ {
		select {
		case res := <-results:
			res.PrintHostOutput()
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	// var wg sync.WaitGroup

	// for _, host := range parsedHosts {
	// 	wg.Add(1)
	// 	rawCmdOpts := cliconfig.CommandOptions{
	// 		Host:               host,
	// 		CommandToRun:       commandToRun,
	// 		Sudo:               false,
	// 		IgnoreHostkeyCheck: ignoreHostkeyCheck,
	// 	}
	// 	go func(host string) {
	// 		defer wg.Done()
	// 		rawCmdOpts.RunCommand()
	// 	}(host)
	// }
	// wg.Wait()
}
