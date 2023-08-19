package cmd

import (
	"fmt"
	"github.com/ajmyyra/bpm/pkg/config"
	"github.com/ajmyyra/bpm/pkg/remote"
	"os"
)

const EnvVarGitHubToken = "GITHUB_TOKEN"

func initGitHubClient() (remote.GitHub, error) {
	ghApiToken, _ := os.LookupEnv(EnvVarGitHubToken)
	gh, err := remote.NewGitHubClient(ghApiToken, fmt.Sprintf("bpm/%s (github.com/ajmyyra/bpm)", config.BPMVersion))
	if err != nil {
		return remote.GitHub{}, err
	}

	return gh, nil
}
