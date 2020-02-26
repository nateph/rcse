package cmd

import (
	"fmt"
	"rcse/pkg/files"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	sequenceFile string
	sequenceCmd  = &cobra.Command{
		Use:   "sequence",
		Short: "Run a sequence of jobs on remote machines.",
		Long:  "Sequence reads from a yaml file and runs the jobs in sequential order, for each machine.",
		Run:   sequenceCommand,
	}
)

func init() {
	cobra.OnInitialize(parseSequenceFile)
	rootCmd.AddCommand(sequenceCmd)
	sequenceCmd.Flags().StringVarP(&sequenceFile, "file", "f", "", "the sequence file, in yaml format.")
	sequenceCmd.MarkFlagRequired("file")
	viper.BindPFlag("sequenceFile", sequenceCmd.Flags().Lookup("file"))
}

func parseSequenceFile() {
	if viper.IsSet("sequenceFile") {
		viper.SetConfigType("yaml")
		sequenceFilePath, err := files.ParseAndVerifyFilePath(sequenceFile)
		if err != nil {
			logrus.Fatal(err.Error())
		}
		viper.SetConfigFile(sequenceFilePath)

		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using sequence file:", viper.ConfigFileUsed())
		}
	}
}

func sequenceCommand(cmd *cobra.Command, args []string) {
	jobName := viper.GetString("shell")
	fmt.Println(jobName)
}
