package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/spf13/cobra"
)

var (
	inventoryFile string
	commandToRun  string
)

var rootCmd = &cobra.Command{
	Use:   "rcse",
	Short: "Run a command somewhere else",
	Run: func(cmd *cobra.Command, args []string) {
		runRootCmd()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rcse.yaml)")
	rootCmd.Flags().StringVarP(&inventoryFile, "inventory", "i", "", "the inventory file of hosts to lookup.")
	rootCmd.Flags().StringVarP(&commandToRun, "command", "c", "", "the command to run on a remote host")
}

func runCommand(host string, command string) {
	cmd := exec.Command("ssh", "-t", "-oStrictHostKeyChecking=no", host, command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed on %s, error was: %s\n", host, err)
		return
	}
	fmt.Printf("%s - %s", host, out.String())
}

func runRootCmd() {
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
