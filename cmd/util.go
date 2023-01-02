package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type RequestedProjectDetails struct {
	Owner   string
	Project string
	Package string
	Version string
}

func parseArgsForProjectDetails(arguments []string) (RequestedProjectDetails, error) {
	repoParts := strings.Split(arguments[0], "/")
	if len(repoParts) != 2 {
		return RequestedProjectDetails{}, errors.New("invalid repository, see example")
	}

	details := RequestedProjectDetails{
		Owner:   repoParts[0],
		Project: strings.Split(repoParts[1], "=")[0],
		Package: repoParts[1],
	}
	if len(arguments) == 2 {
		details.Package = arguments[1]
	}

	if strings.Contains(details.Package, "=") {
		versioned := strings.Split(details.Package, "=")
		if len(versioned) != 2 {
			return RequestedProjectDetails{}, errors.New("erroneously formatted version, see example")
		}
		details.Package = versioned[0]
		details.Version = versioned[1]
	}

	return details, nil
}

func parseDate(time time.Time) string {
	return fmt.Sprintf("%d-%s-%s", time.Year(), int2Char(int(time.Month())), int2Char(time.Day()))
}

func int2Char(num int) string {
	if num < 10 {
		return fmt.Sprintf("0%d", num)
	}
	return fmt.Sprintf("%d", num)
}
