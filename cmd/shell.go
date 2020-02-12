package cmd

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"sync"

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
	cmd := exec.Command("ssh", "-t", "-oStrictHostKeyChecking=no", host, command)
	// var out bytes.Buffer
	// cmd.Stdout = &out
	stdout, err := cmd.Output()
	if err != nil {
		log.Printf("Failed on %s, error was: %s\n", host, err)
		return
	}
	result := common.CommandResult{
		Stdout:     stdout,
		ReturnCode: cmd.ProcessState.ExitCode(),
		Host:       host,
	}
	common.PrintHostOutput(result)
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
