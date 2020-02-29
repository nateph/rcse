package main

import (
	"os"
	"rcse/cmd"
)

func main() {
	rcseCmd := cmd.NewRootCmd(os.Stdout, os.Args[1:])

	if err := rcseCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
