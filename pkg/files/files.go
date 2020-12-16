package files

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// ParseAndVerifyFilePath will take the passed inventory file and parse/expand
// it to an absolute path. It will then check the file exists before returning the path.
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

// VerifyScript checks to make sure a script has a shebang
func VerifyScript(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	fileReader := bufio.NewReader(file)
	shebang, err := fileReader.Peek(2)
	if err != nil {
		return err
	}
	if string(shebang) != "#!" {
		return fmt.Errorf("%s does not have a proper shebang", filePath)
	}

	return nil
}
