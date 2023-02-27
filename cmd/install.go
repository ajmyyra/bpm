package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ajmyyra/bpm/pkg/remote"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a new binary",
	Long: `Install downloads and starts tracking a new binary.
	If a project has more than one binary, you need to specify its name`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")

		if len(args) < 1 || len(args) > 2 {
			fmt.Printf("Invalid amount of arguments.\n%s\n", cmd.Example)
			os.Exit(1)
		}

		gh, err := remote.NewGitHubClient("")
		if err != nil {
			panic(err) // TODO improve
		}

		repoParts := strings.Split(args[0], "/")
		if len(repoParts) != 2 {
			fmt.Println("invalid repository, see example")
			os.Exit(1)
		}

		details, err := gh.GetDetails(repoParts[0], repoParts[1], "")
		if err != nil {
			fmt.Printf("trouble fetching repository details: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("reached end of install, received details: " + details)
	},
	Example: `bpm install kubernetes-sigs/kind
bpm install kubernetes-sigs/cluster-api clusterctl`,
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.PersistentFlags().StringP("source", "s", "GitHub", "Package source")
}
