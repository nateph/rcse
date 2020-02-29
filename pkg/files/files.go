package files

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// InventoryFile should only contain one yaml entry
type InventoryFile struct {
	Hosts []string `yaml:"hosts"`
}

// ParseAndVerifyFilePath will take the passed inventory file string from the flag and
// parse/expand it to an absolute path. It will then check the file exists before returning the path.
func ParseAndVerifyFilePath(filePath string) (string, error) {
	var absFilePath string

	currentUser, _ := user.Current()
	userHomeDir := currentUser.HomeDir
	if strings.HasPrefix(filePath, "~/") {
		absFilePath = filepath.Join(userHomeDir, filePath[2:])
	}
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

	return LoadReader(f)
}

// LoadReader returns the contents of a config file
func LoadReader(fd io.Reader) (inv InventoryFile, err error) {
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return inv, err
	}
	err = yaml.UnmarshalStrict(data, &inv)
	return inv, err
}
