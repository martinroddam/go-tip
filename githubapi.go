package gotip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type GitInfo struct {
	Owner               string `yaml:"Owner"`
	Repo                string `yaml:"Repo"`
	PersonalAccessToken string `yaml:"PersonalAccessToken"`
}

type GitHubAPI []GitHubCommit

type GitHubCommit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
	}
}

const label string = "Verified in PROD"

func getMostRecentlyMergedPullRequest() string {
	var gitInfo = getGitInfo()
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", gitInfo.Owner, gitInfo.Repo)

	client := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Authorization", "token "+gitInfo.PersonalAccessToken)

	res, getErr := client.Do(req)
	if getErr != nil {
		fmt.Println(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Println(readErr)
	}

	var data GitHubAPI
	json.Unmarshal(body, &data)

	r := regexp.MustCompile(`Merge pull request #(?P<prNumber>\d+) .*`)
	var prNumber string

	for i := range data {
		gitHubCommit := data[i].Commit
		if len(r.FindStringSubmatch(gitHubCommit.Message)) > 1 {
			prNumber = r.FindStringSubmatch(gitHubCommit.Message)[1]
			return prNumber
		}
	}
	return ""
}

func applyLabelToPullRequest(prNumber string) {
	var gitInfo = getGitInfo()
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s/labels", gitInfo.Owner, gitInfo.Repo, prNumber)
	client := http.Client{}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader("[\""+label+"\"]"))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Authorization", "token "+gitInfo.PersonalAccessToken)
	req.Header.Set("Content-Type", "application/json")

	res, getErr := client.Do(req)
	if getErr != nil {
		fmt.Println(getErr)
		fmt.Printf("Failed to apply label to PR# [%s]. This will be re-attempted.\n", prNumber)
		unmarkPullRequestVerified(prNumber)
		return
	}
	if res.StatusCode == 200 {
		return
	}
	fmt.Printf("Failed to apply label to PR# [%s]. Github return status [%s]. This will be re-attempted.\n", strconv.Itoa(res.StatusCode), prNumber)
	unmarkPullRequestVerified(prNumber)
}
