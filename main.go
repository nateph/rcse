package main

import (
	"os"

	"github.com/nateph/rcse/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	rcseCmd := cmd.NewRcseCommand(os.Stdout, os.Args[1:])

	if err := rcseCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
