package cliconfig

import (
	"errors"
	"fmt"
	"os"

	"github.com/nateph/rcse/pkg/command"
	"github.com/nateph/rcse/pkg/files"
	"github.com/nateph/rcse/pkg/files/inventory"
	"github.com/spf13/pflag"
)

// Options represents fields under the options key
type Options struct {
	FailureLimit       int
	Forks              int
	IgnoreHostKeyCheck bool
	InventoryFilePath  string
	ListHosts          bool
	OutFormat          string
	Password           string
	PrivateKey         string
	User               string
	Verbose            bool
}

// Job represents a singular job
type Job struct {
	Command string
	Script  string
	// If the job is a script, and Cleanup is true, we delete the script afterwards
	Cleanup bool
}

// Config includes all configuration for the program
type Config struct {
	// We can either be running single command or a script so we abstract it to a job type
	Job     Job
	Options Options
}

// AddBaseFlags binds base flags from the root command to the given flagset.
func (o *Options) AddBaseFlags(fs *pflag.FlagSet) {
	fs.IntVar(&o.FailureLimit, "failure-limit", 1000, "stop execution after n amount of hosts return a failure")
	fs.IntVar(&o.Forks, "forks", 1, "the max amount of hosts to run at any given time")
	fs.BoolVarP(
		&o.IgnoreHostKeyCheck,
		"ignore-hostkey-checking",
		"k",
		false,
		"disable host key verification. this will accept any host key and is insecure\n"+
			"same as 'ssh -o StrictHostKeyChecking=no' ")
	fs.StringVarP(&o.InventoryFilePath, "inventory", "i", o.InventoryFilePath, "the inventory file of hosts to run on, in yaml format")
	fs.BoolVar(&o.ListHosts, "list-hosts", false, "only list hosts that will be executed on, then exits")
	fs.StringVarP(&o.OutFormat, "format", "o", "text", "format result output. takes text/json/yaml")
	fs.StringVarP(&o.Password, "password", "p", o.Password, "the password for a remote user")
	fs.Lookup("password").NoOptDefVal = "default"
	fs.StringVarP(&o.User, "user", "u", o.User, "the optional user to execute as, if -p is not provided, will prompt for password")
	fs.StringVar(&o.PrivateKey, "private-key", o.PrivateKey, "specify a private key to use for ssh connection")
	// fs.BoolVarP(&o.Verbose, "verbose", "v", false, "enable verbose logging")
}

// CheckBaseOptions verifies there was correct options specified
func (o *Options) CheckBaseOptions() error {
	if len(o.InventoryFilePath) == 0 {
		return errors.New("no inventory flag was specified, all rcse operations require an inventory")
	} else if o.ListHosts {
		hosts, err := inventory.LoadInventory(o.InventoryFilePath)
		if err != nil {
			return err
		}
		for _, host := range hosts {
			fmt.Println(host)
		}
		// Stop execution and exit successfully if --list-hosts was passed
		os.Exit(0)
	}

	if o.Forks <= 0 {
		return errors.New("forks value needs to be above 0")
	}

	if len(o.PrivateKey) != 0 {
		_, err := files.ParseAndVerifyFilePath(o.PrivateKey)
		if err != nil {
			return err
		}
	}

	if err := o.VerifyUserOpts(); err != nil {
		return err
	}
	return nil
}

// VerifyUserOpts validates user and password options
func (o *Options) VerifyUserOpts() error {
	var err error
	// Prompt for password if a user was supplied with or without -p, but only if -p was empty
	if (len(o.User) != 0 && o.Password == "default") || (len(o.User) != 0 && len(o.Password) == 0) {
		o.Password, err = command.ConsumePassword(o.User, o.Password)
		if err != nil {
			return err
		}
		// if a password was passed without a user
	} else if len(o.User) == 0 && len(o.Password) != 0 {
		return errors.New("please set a user if supplying a password")
	}
	return nil
}
