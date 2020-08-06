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

// Concurrent holds channel info for workerpool
// type Concurrent struct {
// 	Jobs    chan command.Options
// 	Results chan command.Result
// 	Err     chan error
// }

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
	fs.StringVarP(&o.InventoryFilePath, "inventory", "i", "", "the inventory file of hosts to run on, in yaml format")
	fs.BoolVar(&o.ListHosts, "list-hosts", false, "list hosts that will be ran on. Doesn't execute anything else")
	fs.StringVarP(&o.OutFormat, "format", "o", "text", "format result output. takes text/json/yaml")
	fs.StringVarP(&o.Password, "password", "p", "default", "the password for a remote user supplied by -u or --user")
	fs.Lookup("password").NoOptDefVal = "default"
	fs.StringVarP(&o.User, "user", "u", "", "the optional user to execute as, if -p is not provided, will prompt for password")
}

// CheckBaseOptions verifies there was correct options specified
func (o *Options) CheckBaseOptions() error {
	if o.InventoryFilePath == "" {
		return errors.New("no inventory flag was specified, all rcse operations require an inventory")
	} else if o.ListHosts {
		inventory, err := LoadInventory(o.InventoryFilePath)
		if err != nil {
			return err
		}
		for _, host := range inventory.Hosts {
			fmt.Println(host)
		}
	}

	if o.Forks == 0 {
		return errors.New("forks value needs to be above 0")
	}

	// if --username and --password were supplied correctly without --list-hosts
	var err error
	if o.User != "" && o.Password == "default" && !o.ListHosts {
		o.Password, err = command.ConsumePassword(o.User, o.Password)
		if err != nil {
			return err
		}
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

// LoadInventory returns the inventory file contents
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
