package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	stateFile string
	rootCmd   = &cobra.Command{
		Use:   "bpm",
		Short: "bpm - Binary Package Manager",
		Long:  `Binary Package Manager keeps track of your binary applications and helps to keep them up to date`,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&stateFile, "state", "", "specify a state file location")
}
