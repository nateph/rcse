package main

import (
	"fmt"
	"os"

	"github.com/nateph/rcse/cmd"
)

func main() {
	rcseCmd := cmd.NewRcseCommand(os.Stdout, os.Args[1:])

	if err := rcseCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
