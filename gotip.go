package gotip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

type GitInfo struct {
	Owner               string `yaml:"Owner"`
	Repo                string `yaml:"Repo"`
	PersonalAccessToken string `yaml:"PersonalAccessToken"`
}

type PathInfo struct {
	Paths []Path `yaml:"Paths"`
}

type Path struct {
	PathName string `yaml:"PathName"`
	PathDesc string `yaml:"PathDesc"`
}

var c = cache.New(0, 0) // no expiry
const layout = "2006-01-02T15:04:05.000Z"

func init() {
	//var pathInfo = getPathInfo()
	var gitInfo = getGitInfo()
	go initMostRecentlyMergedPR(gitInfo)
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
		go markGitHubPullRequestAsVerified(lastMergedPR)
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
	paths := getPathInfo().Paths
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
	pathsToBeVerified := getPathInfo().Paths
	for i := range pathsToBeVerified {
		if !isPathVerfiedForThisPullRequest(pathsToBeVerified[i].PathName, prNumber) {
			return false
		}
	}
	fmt.Printf("All paths verified for PR# [%s]\n", prNumber)
	return true
}

func initMostRecentlyMergedPR(gitInfo GitInfo) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", gitInfo.Owner, gitInfo.Repo)
	fmt.Println(url)
	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "token "+gitInfo.PersonalAccessToken)

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	fmt.Println(string(body[:]))

	var data GitHubAPI
	json.Unmarshal(body, &data)
	fmt.Printf("Response: %v\n", data)

	prNumber := "1"
	c.Set("current_pr", prNumber, cache.DefaultExpiration)
	c.Set("PR-"+prNumber, "", cache.DefaultExpiration)
	fmt.Printf("Pull request number: [%s]\n", prNumber)
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
