package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ajmyyra/bpm/pkg/config"
	"github.com/ajmyyra/bpm/pkg/remote"
	"github.com/ajmyyra/bpm/pkg/util"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a new binary",
	Long: `Install downloads a new binary for you and starts tracking it.
	If a project has more than one binary, you need to specify its name`,
	Example: `bpm install kubernetes-sigs/kind
bpm install kubernetes-sigs/cluster-api clusterctl`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 || len(args) > 2 {
			fmt.Println("Invalid amount of arguments. See 'bpm install --help' for examples.")
			os.Exit(1)
		}

		// TODO check for token existence
		// TODO centralize this
		ghApiToken := ""
		gh, err := remote.NewGitHubClient(ghApiToken, fmt.Sprintf("bpm/%s (github.com/ajmyyra/bpm)", config.BPMVersion))
		if err != nil {
			panic(err) // TODO improve
		}

		repoParts := strings.Split(args[0], "/")
		if len(repoParts) != 2 {
			fmt.Println("invalid repository, see example")
			os.Exit(1)
		}

		packageName := repoParts[1]
		if len(args) == 2 {
			packageName = args[1]
		}

		details, err := gh.GetRepositoryDetails(repoParts[0], repoParts[1])
		if err != nil {
			fmt.Printf("trouble fetching repository details: %s\n", err)
			os.Exit(1)
		}

		// TODO move most to bpm show and format it better.

		fmt.Printf("Repo %s, created at %s and updated at %s: %s\n", details.Name, details.CreatedAt, details.UpdatedAt, details.Description)
		fmt.Printf("Fork: %t, Private: %t, releases: %d\n", details.Fork, details.Private, len(details.Releases))

		validName := util.GetEnvMatchingPackage(packageName)
		fmt.Printf("Required package name: %s\n", validName)
		for _, rel := range details.Releases {
			fmt.Printf("Release %s (version %s), published at %s, semversioned: %t\nAssets:\n", rel.Name, rel.Version, rel.PublishedAt, rel.SemVersioned)
			for _, asset := range rel.Assets {
				valid := false
				if validName == asset.Name {
					valid = true
				}
				fmt.Printf("- %s (valid %t), available at %s\n", asset.Name, valid, asset.DownloadURL)
			}

			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.PersistentFlags().StringP("source", "s", "GitHub", "Package source")
}
