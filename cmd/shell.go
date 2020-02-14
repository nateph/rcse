package cmd

import (
	"bufio"
	"log"
	"os"
	"sync"

	"rcse/cmd/cliconfig"
	"rcse/pkg/common"

	"github.com/spf13/cobra"
)

var (
	commandToRun string
	shellCmd     = &cobra.Command{
		Use:   "shell",
		Short: "Execute a shell command.",
		Long:  "Execute a shell command on a remote host.",
		Run:   shellCommand,
	}
)

func init() {
	rootCmd.AddCommand(shellCmd)
	shellCmd.Flags().StringVarP(&commandToRun, "command", "c", "", "the command to run on a remote host")
	shellCmd.MarkFlagRequired("command")
}

func runCommand(host string, command string) {
	// ignoreHostkeyCheck is a persistent flag set in the root command
	sshSession := cliconfig.EstablishSSHConnection(host, ignoreHostkeyCheck)
	defer sshSession.Close()

	stdout, err := sshSession.Output(command)

	if err != nil {
		log.Fatalf("Failed on %s, error was: %s\n", host, err)
		return
	}
	result := common.CommandResult{
		Stdout: stdout,
		// ReturnCode: cmd.ProcessState.ExitCode(),
		Host: host,
	}
	result.PrintHostOutput()
}

func shellCommand(cmd *cobra.Command, args []string) {
	var wg sync.WaitGroup

	hostsFile, err := os.Open(inventoryFile)
	if err != nil {
		log.Fatal(err)
	}
	defer hostsFile.Close()

	scanner := bufio.NewScanner(hostsFile)
	for scanner.Scan() {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			runCommand(host, commandToRun)
		}(scanner.Text())
	}
	wg.Wait()
}
