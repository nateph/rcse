package cmd

import (
	"fmt"
	"io"
	"runtime"

	"github.com/spf13/cobra"
)

type programVersion struct {
	rcse   string
	golang string
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
	version := programVersion{
		rcse:   "0.1",
		golang: runtime.Version(),
	}
	fmt.Printf("rcse v%s\nbuilt by %s", version.rcse, version.golang)
	return nil
}
