package remote

import (
	"context"
	"fmt"
	"golang.org/x/mod/semver"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

const GitHubTime = "2006-01-02 15:04:05 +0000 UTC"

type GitHub struct {
	cli *github.Client
}

func NewGitHubClient(token, useragent string) (GitHub, error) {
	gh := GitHub{}
	var tc *http.Client

	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc = oauth2.NewClient(context.Background(), ts)
	}

	gh.cli = github.NewClient(tc)
	if useragent != "" {
		gh.cli.UserAgent = useragent
	}

	return gh, nil
}

func (s *GitHub) GetRepositoryDetails(owner, repo string) (RepositoryDetails, error) {
	details, _, err := s.cli.Repositories.Get(context.TODO(), owner, repo)
	if err != nil {
		return RepositoryDetails{}, fmt.Errorf("unable to fetch repository details: %w", err)
	}

	createdAt, err := time.Parse(GitHubTime, details.GetCreatedAt().String())
	if err != nil {
		return RepositoryDetails{}, fmt.Errorf("invalid repository creation time: %w", err)
	}

	updatedAt, err := time.Parse(GitHubTime, details.GetUpdatedAt().String())
	if err != nil {
		return RepositoryDetails{}, fmt.Errorf("invalid repository update time: %w", err)
	}

	repoDetails := RepositoryDetails{
		Name:      details.GetFullName(),
		Private:   details.GetPrivate(),
		Fork:      details.GetFork(),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Url:       details.GetHTMLURL(),
	}

	repoDetails.Releases, err = s.GetRepositoryReleases(owner, repo)
	if err != nil {
		return RepositoryDetails{}, err
	}

	return repoDetails, nil
}

func (s *GitHub) GetRepositoryReleases(owner, repo string) ([]RepositoryRelease, error) {
	repoReleases := []RepositoryRelease{}

	releases, _, err := s.cli.Repositories.ListReleases(context.TODO(), owner, repo, &github.ListOptions{})
	if err != nil {
		return repoReleases, fmt.Errorf("unable to fetch repository releases: %w", err)
	}

	for _, rel := range releases {
		version := rel.GetTagName()
		if version == "" {
			version = getVersionFromName(rel.GetName())
		}
		release := RepositoryRelease{
			Name:         rel.GetName(),
			Description:  rel.GetBody(),
			ID:           int(rel.GetID()),
			Draft:        rel.GetDraft(),
			Prerelease:   rel.GetPrerelease(),
			Version:      version,
			SemVersioned: isSemanticallyVersioned(version),
			Url:          rel.GetHTMLURL(),
		}

		publishedAt, err := time.Parse(GitHubTime, rel.PublishedAt.String())
		if err != nil {
			return repoReleases, fmt.Errorf("invalid publication time for release %s: %w", rel.GetName(), err)
		}
		release.PublishedAt = publishedAt

		release.Assets = []ReleaseAsset{}
		for _, GHAsset := range rel.Assets {
			asset := ReleaseAsset{
				ID:          int(GHAsset.GetID()),
				Name:        GHAsset.GetName(),
				DownloadURL: GHAsset.GetBrowserDownloadURL(),
				ContentType: GHAsset.GetContentType(),
			}
			release.Assets = append(release.Assets, asset)
		}

		repoReleases = append(repoReleases, release)
	}

	return repoReleases, nil
}

func getVersionFromName(name string) string {
	parts := strings.Split(name, " ")
	if len(parts) == 0 {
		return ""
	}

	return parts[0]
}

func isSemanticallyVersioned(name string) bool {
	if semver.IsValid(name) {
		return true
	}
	if semver.IsValid(fmt.Sprintf("v%s", name)) {
		return true
	}

	return false
}
