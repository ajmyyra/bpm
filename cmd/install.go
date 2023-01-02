package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a new binary",
	Long: `Install downloads and starts tracking a new binary.
	If a project has more than one binary, you need the specify`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")
	},
	Example: `bpm install kubernetes-sigs/kind
bpm install kubernetes-sigs/cluster-api clusterctl`,
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.PersistentFlags().StringP("source", "s", "GitHub", "Package source")
}
