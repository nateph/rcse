package cliconfig

import "github.com/spf13/pflag"

// CliSettings is a holder for info about the base commands flags
type CliSettings struct {
	IgnoreHostKeyCheck bool
	InventoryFile      string
	ListHosts          bool
	Password           string
	User               string
}

// AddFlags binds base flags from the root command to the given flagset.
func (cs *CliSettings) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&cs.ListHosts, "list-hosts", false, "list hosts that will be ran on. Doesn't execute anything else.")
	fs.StringVarP(&cs.InventoryFile, "inventory", "i", "", "the inventory file of hosts to run on, in yaml format.")
	fs.StringVarP(&cs.User, "user", "u", "", "the optional user to execute as, requires -p")
	fs.StringVarP(&cs.Password, "password", "p", "default", "the password for a remote user supplied by -u or --user.")
	fs.Lookup("password").NoOptDefVal = "default"
	fs.BoolVar(
		&cs.IgnoreHostKeyCheck,
		"ignore-hostkey-checking",
		false,
		"disable host key verification. This will accept any host key and is insecure.\n"+
			"this is the same as 'ssh -o StrictHostKeyChecking=no' ")
}
