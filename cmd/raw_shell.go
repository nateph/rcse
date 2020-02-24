package cmd

import (
	"sync"

	"rcse/cmd/cliconfig"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	commandToRun string
	rawShellCmd  = &cobra.Command{
		Use:   "raw_shell",
		Short: "Execute a shell command",
		Long:  "Execute a shell command on a remote host",
		Run:   rawShellCommand,
	}
)

func init() {
	rootCmd.AddCommand(rawShellCmd)
	rawShellCmd.Flags().StringVarP(&commandToRun, "command", "c", "", "the command to run on a remote host")
	rawShellCmd.MarkFlagRequired("command")
}

func rawShellCommand(cmd *cobra.Command, args []string) {
	parsedHosts := viper.GetStringSlice("hosts")
	var wg sync.WaitGroup

	for _, host := range parsedHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			runCommand(host, commandToRun)
		}(host)
	}
	wg.Wait()
}

func runCommand(host string, command string) {
	// ignoreHostkeyCheck is a persistent flag set in the root command
	sshSession := cliconfig.EstablishSSHConnection(host, ignoreHostkeyCheck)
	defer sshSession.Close()

	result := cliconfig.RunSSHCommand(command, host, sshSession)

	result.PrintHostOutput()
}
