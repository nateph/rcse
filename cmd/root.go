package cmd

import (
	"fmt"
	"os"
	"rcse/cmd/cliconfig"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	listHosts          bool
	inventoryFile      string
	userToBecome       string
	userPassword       string
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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&listHosts, "list-hosts", false, "list hosts that will be ran on. Doesn't execute anything else.")

	rootCmd.PersistentFlags().StringVarP(&inventoryFile, "inventory", "i", "", "the inventory file of hosts to run on, in yaml format.")
	viper.BindPFlag("inventory", rootCmd.PersistentFlags().Lookup("inventory"))

	rootCmd.PersistentFlags().StringVarP(&userToBecome, "user", "u", "", "the optional user to execute as, requires -p")
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))

	rootCmd.PersistentFlags().StringVarP(&userPassword, "password", "p", "default", "the password for a remote user supplied by -u or --user.")
	rootCmd.PersistentFlags().Lookup("password").NoOptDefVal = "default"
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))

	rootCmd.PersistentFlags().BoolVar(
		&ignoreHostkeyCheck,
		"ignore-hostkey-checking",
		false,
		"disable host key verification. This will accept any host key and is insecure.\n"+
			"this is the same as 'ssh -o StrictHostKeyChecking=no' ")

}

// The "inventory file" will be the source of all configuration for the program.
// it can be parsed by each subcommand to have information extracted out as needed
func initConfig() {
	if viper.IsSet("user") && !viper.IsSet("password") {
		logrus.Fatal("user flag was supplied, but a password wasn't.")
	} else if viper.IsSet("password") && viper.GetString("password") == "default" && !listHosts {
		cliconfig.CheckAndConsumePassword()
	}

	if inventoryFile != "" {
		viper.SetConfigType("yaml")
		inventoryFilePath := cliconfig.ParseAndVerifyFile(inventoryFile)
		viper.SetConfigFile(inventoryFilePath)

	}
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
