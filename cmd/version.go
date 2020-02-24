package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version number of rcse",
	Long:  `All software has versions. This is the one for rcse.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rcse v0.1")
	},
}
