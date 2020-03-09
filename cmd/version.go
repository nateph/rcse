package cmd

import (
	"fmt"
	"io"
	"runtime"

	"github.com/spf13/cobra"
)

// Declare to package scope for Make
var (
	buildDate = "unknown"
	gitCommit = "unknown"
	rcse      = "unknown"
)

// ProgramVersion contains info for various versions related to the program
type programVersion struct {
	buildDate string
	gitCommit string
	golang    string
	rcse      string
}

func newVersionCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		Short:                 "Display the version number of rcse",
		Long:                  "All software has versions. This is the one for rcse.",
		RunE:                  runVersion,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

func runVersion(cmd *cobra.Command, args []string) error {
	Version := programVersion{
		gitCommit: gitCommit,
		buildDate: buildDate,
		golang:    runtime.Version(),
		rcse:      rcse,
	}
	fmt.Printf("rcse v%s\n%s\ngit commit %s\nbuilt on %s", Version.rcse, Version.golang, Version.gitCommit, Version.buildDate)
	return nil
}
