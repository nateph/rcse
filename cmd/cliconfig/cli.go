package cliconfig

import "github.com/spf13/pflag"

// CliSettings is a holder for info about the base commands flag values
type CliSettings struct {
	FailureLimit       int
	Forks              int
	IgnoreHostKeyCheck bool
	InventoryFile      string
	ListHosts          bool
	Password           string
	User               string
}

// AddFlags binds base flags from the root command to the given flagset.
func (cs *CliSettings) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&cs.ListHosts, "list-hosts", false, "list hosts that will be ran on. Doesn't execute anything else")
	fs.StringVarP(&cs.InventoryFile, "inventory", "i", "", "the inventory file of hosts to run on, in yaml format")
	fs.StringVarP(&cs.User, "user", "u", "", "the optional user to execute as, requires -p")
	fs.StringVarP(&cs.Password, "password", "p", "default", "the password for a remote user supplied by -u or --user")
	fs.Lookup("password").NoOptDefVal = "default"
	fs.BoolVarP(
		&cs.IgnoreHostKeyCheck,
		"ignore-hostkey-checking",
		"k",
		false,
		"disable host key verification. this will accept any host key and is insecure.\n"+
			"same as 'ssh -o StrictHostKeyChecking=no' ")
	fs.IntVar(&cs.Forks, "forks", 100, "the max amount of hosts to run at any given time")
	fs.IntVar(&cs.FailureLimit, "failure-limit", 1000, "stop execution after n amount of hosts return a failure")
}
