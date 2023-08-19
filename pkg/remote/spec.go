package remote

import (
	"time"
)

type ReleaseAsset struct {
	ID          int
	Name        string
	DownloadURL string
	ContentType string
}

type RepositoryRelease struct {
	Name         string
	Description  string
	Version      string
	SemVersioned bool
	PublishedAt  time.Time
	ID           int
	Draft        bool
	Prerelease   bool
	Assets       []ReleaseAsset
	Url          string
}

type RepositoryDetails struct {
	Name        string
	Description string
	Private     bool
	Fork        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Url         string
	Releases    []RepositoryRelease
}
