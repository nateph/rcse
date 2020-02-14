package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	inventoryFile      string
	ignoreHostkeyCheck bool
	rootCmd            = &cobra.Command{
		Use:   "rcse",
		Short: "Run a command somewhere else",
	}
)

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
	rootCmd.PersistentFlags().StringVarP(&inventoryFile, "inventory", "i", "", "the inventory file of hosts to lookup.")
	rootCmd.MarkPersistentFlagRequired("inventory")

	rootCmd.PersistentFlags().BoolVar(
		&ignoreHostkeyCheck,
		"ignore-hostkey-checking",
		false,
		"disable host key verification. This will accept any host key and is insecure.\n"+
			"this is the same as 'ssh -o StrictHostKeyChecking=no' ")
}
