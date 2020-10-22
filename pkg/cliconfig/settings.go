package cliconfig

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nateph/rcse/pkg/command"
	"github.com/nateph/rcse/pkg/files"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

// InventoryFile should only contain one yaml entry for hosts
type InventoryFile struct {
	Hosts []string `yaml:"hosts"`
}

// Options represents fields under the options key
type Options struct {
	FailureLimit       int    `yaml:"failurelimit,omitempty"`
	Forks              int    `yaml:"forks,omitempty"`
	IgnoreHostKeyCheck bool   `yaml:"insecure,omitempty"`
	InventoryFilePath  string `yaml:"inventory"`
	ListHosts          bool
	OutFormat          string `yaml:"format,omitempty"`
	Password           string `yaml:"password,omitempty"`
	PrivateKey         string `yaml:"privatekey"`
	User               string `yaml:"user,omitempty"`
}

// Job represents a singular job
type Job struct {
	Command string `yaml:"command"`
	Module  string `yaml:"module"`
	Name    string `yaml:"name"`
}

// ShellOptions represents one shell command
type ShellOptions struct {
	Jobs []Job
}

// Config includes all configuration for the program
type Config struct {
	Jobs    []Job   `yaml:"jobs"`
	Options Options `yaml:",omitempty"`
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
		"disable host key verification. this will accept any host key and is insecure.\n"+
			"same as 'ssh -o StrictHostKeyChecking=no' ")
	fs.StringVarP(&o.InventoryFilePath, "inventory", "i", o.InventoryFilePath, "the inventory file of hosts to run on, in yaml format")
	fs.BoolVar(&o.ListHosts, "list-hosts", false, "list hosts that will be ran on. Doesn't execute anything else")
	fs.StringVarP(&o.OutFormat, "format", "o", "text", "format result output. takes text/json/yaml")
	fs.StringVarP(&o.Password, "password", "p", o.Password, "the password for a remote user")
	fs.Lookup("password").NoOptDefVal = "default"
	fs.StringVarP(&o.User, "user", "u", o.User, "the optional user to execute as, if -p is not provided, will prompt for password")
	fs.StringVar(&o.PrivateKey, "private-key", o.PrivateKey, "specify a private key to use for ssh connection.")
}

// CheckBaseOptions verifies there was correct options specified
func (o *Options) CheckBaseOptions() error {
	if len(o.InventoryFilePath) == 0 {
		return errors.New("no inventory flag was specified, all rcse operations require an inventory")
	} else if o.ListHosts {
		inventory, err := LoadInventory(o.InventoryFilePath)
		if err != nil {
			return err
		}
		for _, host := range inventory.Hosts {
			fmt.Println(host)
		}
		// Stop execution if --list-hosts was passed
		os.Exit(0)
	}

	if o.Forks == 0 {
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

// LoadFile returns the contents of a file as a byte slice
func LoadFile(file string) (data []byte, err error) {
	absFilePath, err := files.ParseAndVerifyFilePath(file)
	if err != nil {
		return data, err
	}
	f, err := os.Open(absFilePath)
	if err != nil {
		return data, err
	}
	defer f.Close()

	data, err = ioutil.ReadAll(f)
	if err != nil {
		return data, err
	}
	return data, nil
}

// LoadInventory returns the inventory file contents as an InventoryFile
func LoadInventory(file string) (inv InventoryFile, err error) {
	data, err := LoadFile(file)
	if err != nil {
		return inv, err
	}
	err = yaml.UnmarshalStrict(data, &inv)
	if err != nil {
		return InventoryFile{}, err
	}
	return inv, nil
}

// LoadConfig reads in a sequence yaml file and stores its information
func LoadConfig(file string) (config *Config, err error) {
	data, err := LoadFile(file)
	if err != nil {
		return config, err
	}
	err = yaml.UnmarshalStrict(data, &config)
	if err != nil {
		return &Config{}, err
	}

	return config, nil
}
