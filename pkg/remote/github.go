package remote

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v50/github"
)

type GitHub struct {
	cli *github.Client
}

func NewGitHubClient(token string) (GitHub, error) {
	gh := GitHub{}
	if token == "" {
		gh.cli = github.NewClient(nil)
	} else {
		return gh, errors.New("using GH tokens is not yet supported")
	}
	gh.cli.UserAgent = "bpm/0.0.1 (github.com/ajmyyra/bpm)"

	return gh, nil
}

func (s *GitHub) GetDetails(owner, repo, version string) (string, error) {
	releases, _, err := s.cli.Repositories.ListReleases(context.Background(), owner, repo, &github.ListOptions{})
	if err != nil {
		panic(err) // at the disco! also change interface & return the error instead
	}

	// TODO check version in releases for matches

	for _, rel := range releases {
		fmt.Printf("Release %s, published at %s, asset url: %s, assets: %+v\n", *rel.Name, rel.PublishedAt, rel.GetAssetsURL(), rel.Assets)
		fmt.Println("")
	}

	return "", nil
}
