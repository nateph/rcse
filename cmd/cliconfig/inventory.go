package cliconfig

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// ParseAndVerifyFile will take the passed inventory file string from the flag and
// parse/expand it to an absolute path.
// It will then check it exists before returning the path.
func ParseAndVerifyFile(filePath string) string {
	var absFilePath string

	currentUser, _ := user.Current()
	userHomeDir := currentUser.HomeDir
	if strings.HasPrefix(filePath, "~/") {
		absFilePath = filepath.Join(userHomeDir, filePath[2:])
	}
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		logrus.Fatalf("Couldn't parse filepath to absolute: %s", filePath)
	}

	fileInfo, err := os.Stat(absFilePath)
	if err != nil || fileInfo.IsDir() {
		logrus.Fatalf("File does not exist: %s", absFilePath)
	}

	return absFilePath
}
