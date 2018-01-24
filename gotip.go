package gotip

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

// PathInfo contains the list of critical path (or journey) labels that
// require completion for the pull request to be considered 'Verified'
type PathInfo struct {
	Paths []Path `yaml:"Paths"`
}

// Path represents a critical path (or journey) in the application.
// It consists of a unique label (Path) and a description. The description
// serves no purpose other than to provide context to the reader.
type Path struct {
	PathName string `yaml:"PathName"`
	PathDesc string `yaml:"PathDesc"`
}

var c = cache.New(0, 0) // go-cache no expiry
const layout = "2006-01-02T15:04:05.000Z"

func init() {
	go initMostRecentlyMergedPR()
}

// Verify a critical path in your application. Must match a PathName as defined in paths.yaml
func Verify(pathName string) {

	lastMergedPR := getMostRecentlyMergedPullRequestFromCache()

	if lastMergedPR == "" {
		fmt.Printf("Pull request number not set. Cannot proceed with verification of path [%s].\n", pathName)
		return
	}

	if !isValidPath(pathName) {
		return
	}

	if isPullRequestAlreadyVerified(lastMergedPR) {
		return
	}

	if isPathVerfiedForThisPullRequest(pathName, lastMergedPR) {
		return
	}

	markPathVerified(pathName, lastMergedPR)

	if areAllPathsVerifiedForPR(lastMergedPR) {
		markPullRequestVerified(lastMergedPR)
		go applyLabelToPullRequest(lastMergedPR)
	}
}

func isValidPath(pathNameToValidate string) bool {
	pathsToBeVerified, err := getPathsToBeVerified()
	if err != nil {
		fmt.Println(err)
		return false
	}
	for i := range pathsToBeVerified {
		if pathsToBeVerified[i].PathName == pathNameToValidate {
			return true
		}
	}
	fmt.Printf("Path Name specified [%s] is not in the configured list\n", pathNameToValidate)
	return false
}

func getPathsToBeVerified() ([]Path, error) {
	paths := getPathInfo().Paths
	if len(paths) < 1 {
		return nil, errors.New("no paths specified")
	}
	return paths, nil
}

func isPathVerfiedForThisPullRequest(pathToBeVerified string, lastMergedPR string) bool {
	pr, found := c.Get(pathToBeVerified)
	if found {
		if strings.Compare(pr.(string), "PR-"+lastMergedPR) == 0 {
			return true
		}
		return false
	}
	return false
}

func isPullRequestAlreadyVerified(mostRecentPullRequestNumber string) bool {
	verifiedTimestamp, found := c.Get("PR-" + mostRecentPullRequestNumber)
	if found {
		if strings.Compare(verifiedTimestamp.(string), "") == 0 {
			return false
		}

		_, err := time.Parse(layout, verifiedTimestamp.(string))

		if err != nil {
			fmt.Println(err)
		}
		return true
	}
	return false
}

func areAllPathsVerifiedForPR(prNumber string) bool {
	pathsToBeVerified := getPathInfo().Paths
	for i := range pathsToBeVerified {
		if !isPathVerfiedForThisPullRequest(pathsToBeVerified[i].PathName, prNumber) {
			return false
		}
	}
	return true
}

func initMostRecentlyMergedPR() {
	prNumber := getMostRecentlyMergedPullRequest()
	c.Set("current_pr", prNumber, cache.DefaultExpiration)
	c.Set("PR-"+prNumber, "", cache.DefaultExpiration)
}

func getMostRecentlyMergedPullRequestFromCache() string {
	pr, found := c.Get("current_pr")
	if found {
		return pr.(string)
	}
	return ""
}

func markPathVerified(pathName string, prNumber string) {
	c.Set(pathName, "PR-"+prNumber, cache.DefaultExpiration)
}

func markPullRequestVerified(prNumber string) {
	t := time.Now()
	timestamp := t.Format(layout)
	c.Set("PR-"+prNumber, timestamp, cache.DefaultExpiration)
}

func unmarkPullRequestVerified(prNumber string) {
	c.Set("PR-"+prNumber, "", cache.DefaultExpiration)
	fmt.Printf("PR# [%s] verifation reset for PR# [%s]\n", prNumber)
}
