package cliconfig

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

// ReadInventoryFile is a wrapper command around viper to read in hosts that we
// want to execute jobs on
func ReadInventoryFile() []string {
	invFilePath, invFileName := parseFilePath()
	invFile := viper.New()
	invFile.SetConfigType("yaml")
	invFile.SetConfigName(invFileName)
	invFile.AddConfigPath(invFilePath)

	err := invFile.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("could not parse inventory file: %s", err))
	}

	return invFile.GetStringSlice("hosts")
}

// parseFilePath will read the inventory flag bound to viper and determine the
// file path and the file name
func parseFilePath() (string, string) {
	var filePath string
	var fileName string

	if viper.IsSet("inventory") {
		passedFile, err := filepath.Abs(viper.GetString("inventory"))
		if err != nil {
			log.Fatalf("Couldn't parse filepath to absolute: %s", viper.GetString("inventory"))
		}
		filePath = filepath.Dir(passedFile)
		fileName = filepath.Base(passedFile)
	} else {
		log.Fatal("Viper could not find bound flag 'inventory'. Exiting")
	}
	return filePath, fileName
}
