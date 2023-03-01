package util

import (
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/mod/semver"
)

func GetVersionFromName(name string) string {
	parts := strings.Split(name, " ")
	if len(parts) == 0 {
		return ""
	}

	return parts[0]
}

func IsSemanticallyVersioned(name string) bool {
	if semver.IsValid(name) {
		return true
	}
	if semver.IsValid(fmt.Sprintf("v%s", name)) {
		return true
	}

	return false
}

func GetEnvMatchingPackage(name string) string {
	return fmt.Sprintf("%s-%s-%s", name, runtime.GOOS, runtime.GOARCH)
}
