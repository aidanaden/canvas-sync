package credit

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/google/go-github/v55/github"
)

// ContributorsSince use GitHub APIs to get all commit since `since` and extract the author information from that. Then
// it will count the author commit and return a sorted list of authors according to number of commits (descending).
func ContributorsSince(since time.Time) []string {
	// Get list of contributors from GitHub APIs.
	ghc := github.NewClient(nil)
	commits, _, err := ghc.Repositories.ListCommits(
		context.Background(),
		"aidanaden",
		"canvas-sync",
		&github.CommitsListOptions{Since: since},
	)

	// If any error happened, print it and exit with code 1.
	if err != nil {
		fmt.Printf("err: %+v", err)
		os.Exit(1)
	}

	// Generate a unique list of contributors.
	contributors := make(map[string]int, 10)
	for _, c := range commits {
		username := c.GetAuthor().GetLogin()
		commitCount, _ := contributors[username]
		contributors[username] = commitCount + 1
	}

	// Sort contributors based on the number of commits (descending).
	sortedContributors := make([]string, 0, len(contributors))
    for username := range contributors {
        sortedContributors = append(sortedContributors, username)
    }
	sort.Slice(sortedContributors, func(i, j int) bool {
		return contributors[sortedContributors[i]] > contributors[sortedContributors[j]]
	})

	// Format the contributors to a list of string.
	result := make([]string, 0, len(contributors))
	for _, username := range sortedContributors {
        commits := "commits"
        if contributors[username] == 1 {
            commits = "commit"
        }

		result = append(result, fmt.Sprintf("%s (%d %s)", username, contributors[username], commits))
	}

	return result
}

// Contributors use GitHub APIs to get 500 top contributors. The list is sorted base on the number of commits
// (descending).
func Contributors() []string {
	// Get list of contributors from GitHub APIs.
	ghc := github.NewClient(nil)
	contributors, _, err := ghc.Repositories.ListContributors(
		context.Background(),
		"aidanaden",
		"canvas-sync",
		nil,
	)

	// If any error happened, print it and exit with code 1.
	if err != nil {
		fmt.Printf("err: %+v", err)
		os.Exit(1)
	}

	// Format the contributors to a list of string.
	result := make([]string, 0, len(contributors))
	for _, c := range contributors {
		if c == nil {
			continue
		}

        commits := "commits"
        if c.GetContributions() == 1 {
            commits = "commit"
        }

		result = append(result, fmt.Sprintf("%s (%d %s)", c.GetLogin(), c.GetContributions(), commits))
	}

	return result
}
