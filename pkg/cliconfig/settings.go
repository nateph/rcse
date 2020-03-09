package cliconfig

import (
	"errors"

	"github.com/spf13/pflag"
)

// BaseSettings is a holder for info about the base commands flag values
type BaseSettings struct {
	FailureLimit       int
	Forks              int
	IgnoreHostKeyCheck bool
	InventoryFile      string
	ListHosts          bool
	OutFormat          string
	Password           string
	User               string
}

// AddFlags binds base flags from the root command to the given flagset.
func (bs *BaseSettings) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&bs.FailureLimit, "failure-limit", 1000, "stop execution after n amount of hosts return a failure")
	fs.IntVar(&bs.Forks, "forks", 100, "the max amount of hosts to run at any given time")
	fs.BoolVarP(
		&bs.IgnoreHostKeyCheck,
		"ignore-hostkey-checking",
		"k",
		false,
		"disable host key verification. this will accept any host key and is insecure.\n"+
			"same as 'ssh -o StrictHostKeyChecking=no' ")
	fs.StringVarP(&bs.InventoryFile, "inventory", "i", "", "the inventory file of hosts to run on, in yaml format")
	fs.BoolVar(&bs.ListHosts, "list-hosts", false, "list hosts that will be ran on. Doesn't execute anything else")
	fs.StringVarP(&bs.OutFormat, "format", "o", "text", "format result output. takes text/json/yaml")
	fs.StringVarP(&bs.Password, "password", "p", "default", "the password for a remote user supplied by -u or --user")
	fs.Lookup("password").NoOptDefVal = "default"
	fs.StringVarP(&bs.User, "user", "u", "", "the optional user to execute as, requires -p")
}

// CheckBaseOptions verifies there was correct options specified
func CheckBaseOptions(bs *BaseSettings) error {
	if bs.InventoryFile == "" {
		return errors.New("no inventory flag was specified, all rcse operations require an inventory")
	}
	if bs.Forks == 0 {
		return errors.New("forks value needs to be above 0")
	}
	return nil
}

// JobOptions contains info for a single job
type JobOptions struct {
	CommandToRun string
	*BaseSettings
}
