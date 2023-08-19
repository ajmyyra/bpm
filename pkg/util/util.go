package util

import (
	"fmt"
	"github.com/ajmyyra/bpm/pkg/remote"
	"runtime"
	"strings"
)

func GetEnvMatchingPackageNames(name, version string) []string {
	var archAltNames = map[string][]string{
		"linux/386":   {"linux32"},
		"linux/amd64": {"linux64"},
		"linux/arm64": {"linux-arm"},
	}

	names := []string{
		fmt.Sprintf("%s-%s-%s", name, runtime.GOOS, runtime.GOARCH),
		fmt.Sprintf("%s-%s-%s-%s", name, version, runtime.GOOS, runtime.GOARCH),
	}

	for _, archName := range archAltNames[fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)] {
		names = append(
			names,
			fmt.Sprintf("%s-%s", name, archName),
			fmt.Sprintf("%s-%s-%s", name, version, archName),
		)
	}

	return names
}

func FindMatchingReleaseAsset(name, version string, assets []remote.ReleaseAsset) *remote.ReleaseAsset {
	possibleNames := GetEnvMatchingPackageNames(name, version)
	for _, asset := range assets {
		for _, name := range possibleNames {
			if strings.ToLower(asset.Name) == name {
				return &asset
			}
		}
	}

	return nil
}
