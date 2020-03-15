package files

import (
	"fmt"
	"os"
	"path/filepath"
)

// InventoryFile should only contain one yaml entry for hosts
// type InventoryFile struct {
// 	Hosts []string `yaml:"hosts"`
// }

// // Options represents fields under the options key
// type Options struct {
// 	FailureLimit       int    `yaml:"failurelimit,omitempty"`
// 	Forks              int    `yaml:"forks,omitempty"`
// 	IgnoreHostKeyCheck bool   `yaml:"insecure,omitempty"`
// 	OutFormat          string `yaml:"format,omitempty"`
// 	Password           string `yaml:"password,omitempty"`
// 	User               string `yaml:"user,omitempty"`
// }

// // Job represents a singular job
// type Job struct {
// 	Command string `yaml:"command"`
// 	Module  string `yaml:"module"`
// 	Name    string `yaml:"name"`
// }

// // Config includes all configuration for the program
// type Config struct {
// 	InvFile InventoryFile
// 	Jobs    []Job   `yaml:"jobs"`
// 	Options Options `yaml:",omitempty"`
// }

// ParseAndVerifyFilePath will take the passed inventory file string from the flag and
// parse/expand it to an absolute path. It will then check the file exists before returning the path.
func ParseAndVerifyFilePath(filePath string) (string, error) {
	var absFilePath string

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("Couldn't parse filepath to absolute: %s", filePath)
	}

	fileInfo, err := os.Stat(absFilePath)
	if os.IsNotExist(err) || fileInfo.IsDir() {
		return "", fmt.Errorf("File does not exist (or is a directory): %s", absFilePath)
	}

	return absFilePath, nil
}

// // LoadInventory returns the inventory file contents
// func LoadInventory(file string) (inv InventoryFile, err error) {
// 	absFilePath, err := ParseAndVerifyFilePath(file)
// 	if err != nil {
// 		return inv, err
// 	}
// 	f, err := os.Open(absFilePath)
// 	if err != nil {
// 		return inv, err
// 	}
// 	defer f.Close()

// 	inventory, err := ioutil.ReadAll(f)
// 	if err != nil {
// 		return inv, err
// 	}
// 	err = yaml.UnmarshalStrict(inventory, &inv)

// 	return inv, nil
// }

// // LoadReader returns the contents of a config file as a Project
// func LoadReader(fd io.Reader) (config Config, err error) {
// 	data, err := ioutil.ReadAll(fd)
// 	if err != nil {
// 		return config, err
// 	}
// 	err = yaml.UnmarshalStrict(data, &config)
// 	return config, err
// }

// // LoadConfig reads in a sequence yaml file and stores its information
// func LoadConfig(file string, invFile string) (config *Config, err error) {
// 	absFilePath, err := ParseAndVerifyFilePath(file)
// 	if err != nil {
// 		return config, err
// 	}
// 	f, err := os.Open(absFilePath)
// 	if err != nil {
// 		return config, err
// 	}
// 	defer f.Close()

// 	data, err := ioutil.ReadAll(f)
// 	if err != nil {
// 		return config, err
// 	}
// 	err = yaml.UnmarshalStrict(data, &config)

// 	inv, err := LoadInventory(invFile)
// 	if err != nil {
// 		return config, err
// 	}
// 	config.InvFile = inv

// 	return config, nil
// }
