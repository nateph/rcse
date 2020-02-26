package files

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

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