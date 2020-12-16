package inventory

import (
	"bufio"
	"os"

	"github.com/nateph/rcse/pkg/files"
	"github.com/nateph/rcse/pkg/files/hostglob"
)

// LoadFile returns the file after it has been verified to exist
func LoadFile(file string) (f *os.File, err error) {
	absFilePath, err := files.ParseAndVerifyFilePath(file)
	if err != nil {
		return nil, err
	}
	verfiedFile, err := os.Open(absFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return verfiedFile, nil
}

// LoadInventory returns the inventory file contents as an InventoryFile
func LoadInventory(file string) (inv []string, err error) {
	f, err := LoadFile(file)
	if err != nil {
		return inv, err
	}

	var hosts []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if hostglob.IsGlob(scanner.Text()) {
			uncollapsed, err := hostglob.Uncollapse(scanner.Text())
			if err != nil {
				return []string{}, nil
			}
			hosts = append(hosts, uncollapsed...)
		} else {
			hosts = append(hosts, scanner.Text())
		}
	}

	return hosts, nil
}