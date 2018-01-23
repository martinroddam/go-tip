package gotip

import (
	"fmt"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

type Config struct {
	RepoURL string `yaml:"RepoURL"`
	Paths   []Path `yaml:"Paths"`
}

type Path struct {
	PathName string `yaml:"PathName"`
	PathDesc string `yaml:"PathDesc"`
}

// set up the go-cache defaults
var c = cache.New(0, 0) // no expiry
//const githubRoot = "api.github.com"
//const gitOwner = "utilitywarehouse"
const layout = "2006-01-02T15:04:05.000Z"

func init() {

	var config = getConfig()

	prNumber := initMostRecentlyMergedPR(config.RepoURL)
	c.Set("current_pr", prNumber, cache.DefaultExpiration)
	c.Set("PR-"+prNumber, "", cache.DefaultExpiration)
	mostRecentPR := getMostRecentlyMergedPR()
	fmt.Println(mostRecentPR)

}

// Verify a path in your application. Must match the PathName as defined in paths.yaml
func Verify(pathName string) {

	lastMergedPR := getMostRecentlyMergedPR()

	isValidPath(pathName)

	if isPullRequestAlreadyVerified(lastMergedPR) {
		return
	}

	if isPathVerfiedForThisPullRequest(pathName, lastMergedPR) {
		return
	}

	markPathVerified(pathName, lastMergedPR)

	if areAllPathsVerifiedForPR(lastMergedPR) {
		markPullRequestVerified(lastMergedPR)
		markGitHubPullRequestAsVerified(lastMergedPR)
	}
}

func isValidPath(pathNameToValidate string) bool {
	pathsToBeVerified, err := getPathsToBeVerified()
	if err != nil {
		fmt.Println(err)
	}
	for i := range pathsToBeVerified {
		if pathsToBeVerified[i].PathName == pathNameToValidate {
			fmt.Printf("Path Name [%s] is valid.\n", pathNameToValidate)
			return true
		}
	}
	fmt.Printf("Path Name [%s] is invalid!\n", pathNameToValidate)
	return false
}

func getPathsToBeVerified() ([]Path, error) {
	paths := getConfig().Paths
	return paths, nil
}

func getPathsAlreadyVerifiedForThisPullRequest(prNumber string) []string {
	// using the list of paths to be verified, which ones have a gp-cache key with the value matching the PR ID
	paths := []string{"one", "two"}
	return paths
}

func isPathVerfiedForThisPullRequest(pathToBeVerified string, lastMergedPR string) bool {
	pr, found := c.Get(pathToBeVerified)
	if found {
		if strings.Compare(pr.(string), "PR-"+lastMergedPR) == 0 {
			fmt.Printf("Path [%s] has been previously verified for PR# [%s]\n", pathToBeVerified, lastMergedPR)
			return true
		}
		fmt.Printf("Path [%s] has not yet been verified for PR# [%s]\n", pathToBeVerified, lastMergedPR)
		return false
	}
	fmt.Printf("Path [%s] has not yet been verified for PR# [%s]\n", pathToBeVerified, lastMergedPR)
	return false
}

func isPullRequestAlreadyVerified(mostRecentPullRequestNumber string) bool {
	verifiedTimestamp, found := c.Get("PR-" + mostRecentPullRequestNumber)
	if found {
		if strings.Compare(verifiedTimestamp.(string), "") == 0 {
			fmt.Println("Pull request not yet verified.")
			return false
		}
		layout := "2006-01-02T15:04:05.000Z"

		_, err := time.Parse(layout, verifiedTimestamp.(string))

		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Pull request already verified at %s.\n", verifiedTimestamp.(string))
		return true

	}
	return false
}

func areAllPathsVerifiedForPR(prNumber string) bool {
	pathsToBeVerified := getConfig().Paths
	for i := range pathsToBeVerified {
		if !isPathVerfiedForThisPullRequest(pathsToBeVerified[i].PathName, prNumber) {
			return false
		}
	}
	fmt.Printf("All paths verified for PR# [%s]\n", prNumber)
	return true
}

func initMostRecentlyMergedPR(repoURL string) string {
	// get this from GitHub
	//url = "https://%s/%s/%s/"
	// return a string of the PR ID only and store this in a key value pair as key: {PR_ID}, value: nil
	prNumber := "1"
	fmt.Printf("Pull request number: [%s]\n", prNumber)
	return prNumber
}

func getMostRecentlyMergedPR() string {
	pr, found := c.Get("current_pr")
	if found {
		fmt.Printf("Current Pull Request ID: [%s]\n", pr.(string))
		return pr.(string)
	}
	return "" // if we can't find it, we should log out but not panic
}

func markPathVerified(pathName string, prNumber string) {
	c.Set(pathName, "PR-"+prNumber, cache.DefaultExpiration)
	fmt.Printf("Path [%s] marked verified for PR# [%s]\n", pathName, prNumber)
}

func markPullRequestVerified(prNumber string) {
	t := time.Now()
	timestamp := t.Format(layout)
	c.Set("PR-"+prNumber, timestamp, cache.DefaultExpiration)
	fmt.Printf("PR# [%s] marked verified at [%s]\n", prNumber, timestamp)
}

func markGitHubPullRequestAsVerified(prNumber string) {
	fmt.Printf("Pull request # [%s] tagged as Verified in GitHub.", prNumber)
}
