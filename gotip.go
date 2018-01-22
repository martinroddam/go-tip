package gotip

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/patrickmn/go-cache"
	"gopkg.in/yaml.v2"
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

func init() {

	var config Config
	filename := "./paths.yaml"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	ymlErr := yaml.Unmarshal(content, &config)
	if ymlErr != nil {
		log.Fatalf("error: %v", ymlErr)
	}

	prNumber := initMostRecentlyMergedPR(config.RepoURL)
	c.Set("current_pr", prNumber, cache.DefaultExpiration)
	c.Set(prNumber, "", cache.DefaultExpiration)

}

func verify(path string) {
	// do something
}

func getPathsToBeVerified() ([]string, error) {
	return nil, nil
}

func pathExists(path string) bool {
	// load the paths from the config and check for existence
	return true
}

func getLastMergedPullRequestNumber(githubProject string) string {
	// load the project id from the config yaml
	// get the PR ID most recently merged
	return ""
}

func isPullRequestAlreadyVerified(prNumber string) bool {
	// look in go-cache for a key matching the PR ID
	// value nil implies not verified - a verified PR has a time stamp
	return false
}

func getPathsAlreadyVerifiedForThisPullRequest(prNumber string) []string {
	// using the list of paths to be verified, which ones have a gp-cache key with the value matching the PR ID
	paths := []string{"one", "two"}
	return paths
}

func isPathVerfiedForThisPullRequest(path string, prNumber string) bool {
	// is there a key/value pair in go-cache matching this path and PR ID
	return false
}

func areAllPathsVerified(prNumber string) bool {
	// get all paths from go-cache and compare to the list of paths from the config
	return true
}

func initMostRecentlyMergedPR(repoUrl string) string {
	// get this from GitHub

	// return a string of the PR ID only and store this in a key value pair as key: {PR_ID}, value: nil
	return "179"
}

func getMostRecentlyMergedPR() string {
	pr, found := c.Get("current_pr")
	if found {
		fmt.Printf("Current Pull Request ID: [%s]\n", pr.(string))
		return pr.(string)
	}
	return "" // if we can't find it, we should log out but not panic

}

func markPullRequestVerified(prNumber string) {
	// create an entry in go-cache with key: {PR_ID}, value: timestamp
}

func markGitHubPullRequestAsVerified(prNumber string) {
	// tag the PR as verified in Production
}
