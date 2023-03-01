package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall local installation of a package",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("uninstall called")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
