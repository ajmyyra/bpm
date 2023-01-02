package cmd

import (
	"fmt"
	"github.com/ajmyyra/bpm/pkg/util"
	"os"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show binary details and available versions",
	Long:  `Show gives you binary details, currently installed version and other available versions.`,
	Example: `bpm show kubernetes-sigs/kind
bpm show kubernetes-sigs/cluster-api clusterctl`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 || len(args) > 2 {
			fmt.Println("Invalid amount of arguments. See 'bpm install --help' for examples.")
			os.Exit(1)
		}

		gh, err := initGitHubClient()
		if err != nil {
			fmt.Printf("GitHub client initialization failed: %w\n", err)
			os.Exit(2)
		}

		projDetails, err := parseArgsForProjectDetails(args)
		if err != nil {
			fmt.Printf("%w\n", err)
			os.Exit(1)
		}

		details, err := gh.GetRepositoryDetails(projDetails.Owner, projDetails.Project)
		if err != nil {
			fmt.Printf("trouble fetching repository details: %s\n", err)
			os.Exit(1)
		}

		fmt.Println(details.Name)
		fmt.Printf("  Created: %s\n", parseDate(details.CreatedAt))
		fmt.Printf("  Updated: %s\n", parseDate(details.UpdatedAt))
		fmt.Println("  Releases:")
		// TODO have --detailed flag for more details like Description

		if len(details.Releases) == 0 {
			fmt.Println("    No releases yet")
		}

		// TODO only show 10 latest versions. if installed package is somewhere further,
		// add some "... (9 releases omitted)" (unless --detailed) before and show that and "..." below (unless --detailed)
		for _, rel := range details.Releases {
			asset := util.FindMatchingReleaseAsset(projDetails.Package, rel.Version, rel.Assets)
			assetMsg := fmt.Sprintf("no suitable asset available, see %s", rel.Url)
			if asset != nil {
				assetMsg = asset.DownloadURL
			}

			fmt.Printf("    %s (%s) %s\n", rel.Version, parseDate(rel.PublishedAt), assetMsg)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
