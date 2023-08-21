package cmd

import (
	"fmt"
	"github.com/ajmyyra/bpm/pkg/util"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a new binary",
	Long: `Install downloads a new binary for you and starts tracking it.
	If a project has more than one binary, you need to specify its name`,
	Example: `bpm install kubernetes-sigs/kind
bpm install justjanne/powerline-go=v1.24
bpm install kubernetes-sigs/cluster-api clusterctl`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 || len(args) > 2 {
			fmt.Println("Invalid amount of arguments. See 'bpm install --help' for examples.")
			os.Exit(1)
		}

		gh, err := initGitHubClient()
		if err != nil {
			fmt.Printf("GitHub client initialization failed: %w\n", err)
			os.Exit(1)
		}

		projDetails, err := parseArgsForProjectDetails(args)
		if err != nil {
			fmt.Printf("%w\n", err)
			os.Exit(1)
		}

		details, err := gh.GetRepositoryDetails(projDetails.Owner, projDetails.Project)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(details.Releases) == 0 {
			fmt.Printf("No releases available, see %s\n", details.Url)
			os.Exit(1)
		}

		release := details.Releases[0] // GitHub returns the newest version first
		if projDetails.Version != "" {
			match := false
			versions := []string{}
			for _, rel := range details.Releases {
				versions = append(versions, rel.Version)
				if rel.Version == projDetails.Version {
					release = rel
					match = true
					break
				}
			}

			if !match {
				fmt.Printf("Version %s does not match any release.\nAvailable versions: %s\n", projDetails.Version, strings.Join(versions, ", "))
				os.Exit(1)
			}
		}

		asset := util.FindMatchingReleaseAsset(projDetails.Package, release.Version, release.Assets)
		if asset == nil {
			fmt.Printf("No suitable asset available, see %s\n", release.Url)
			os.Exit(1)
		}

		// TODO actual location from config
		if err = util.DownloadAndMoveToLocation(asset.DownloadURL, "/tmp/", projDetails.Package); err != nil {
			fmt.Printf("Downloading asset %s failed: %w\n", asset.DownloadURL, err)
			os.Exit(1)
		}

		fmt.Printf("Package %s (version %s) downloaded succesfully.\n", projDetails.Package, release.Version)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.PersistentFlags().StringP("source", "s", "GitHub", "Package source")
}
