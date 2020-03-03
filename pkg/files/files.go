package files

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// InventoryFile should only contain one yaml entry for hosts
type InventoryFile struct {
	Hosts []string `yaml:"hosts"`
}

// Project includes all project configuration
type Project struct {
	InvFile InventoryFile `yaml:",inline"`
}

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

// LoadInventory returns the inventory file contents
func LoadInventory(file string) (inv InventoryFile, err error) {
	absFilePath, err := ParseAndVerifyFilePath(file)
	if err != nil {
		return inv, err
	}
	f, err := os.Open(absFilePath)
	if err != nil {
		return inv, err
	}
	defer f.Close()

	config, err := LoadReader(f)
	if err != nil {
		return inv, err
	}

	return config.InvFile, nil
}

// LoadReader returns the contents of a config file as a Project
func LoadReader(fd io.Reader) (config Project, err error) {
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return config, err
	}
	err = yaml.UnmarshalStrict(data, &config)
	return config, err
}
