package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

type programVersion struct {
	rcse   string
	golang string
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version number of rcse",
	Long:  `All software has versions. This is the one for rcse.`,
	Run: func(cmd *cobra.Command, args []string) {
		version := programVersion{
			rcse:   "0.1",
			golang: runtime.Version(),
		}
		fmt.Printf("rcse v%s\nbuilt by %s", version.rcse, version.golang)
	},
}
