package gotip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type GitHubAPI []GitHubCommit

type GitHubCommit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
	}
}

type GitHubEditInfo struct {
	labels []string
}

const label string = "Verified in PROD"

func getMostRecentlyMergedPullRequest(gitInfo GitInfo) string {
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

	var data GitHubAPI
	json.Unmarshal(body, &data)

	r := regexp.MustCompile(`Merge pull request #(?P<prNumber>\d+) .*`)
	var prNumber string

	for i := range data {
		gitHubCommit := data[i].Commit
		if len(r.FindStringSubmatch(gitHubCommit.Message)) > 1 {
			prNumber = r.FindStringSubmatch(gitHubCommit.Message)[1]
			fmt.Printf("Found Pull Request Number: %s\n", prNumber)
			return prNumber
		}
	}
	return ""
}

func applyLabelToPullRequest(prNumber string) {
	var gitInfo = getGitInfo()
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s/labels", gitInfo.Owner, gitInfo.Repo, prNumber)
	fmt.Println(url)
	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	//var gitHubEditInfo GitHubEditInfo
	//gitHubEditInfo.labels = append(gitHubEditInfo.labels, label)
	//reqBody, jsonErr := json.Marshal(gitHubEditInfo)
	//if jsonErr != nil {
	//	log.Fatal(jsonErr)
	//}
	//fmt.Println("Applying label: \n" + string(reqBody[:]))
	//req, err := http.NewRequest(http.MethodPatch, url, strings.NewReader(string(reqBody[:])))
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader("[\""+label+"\"]"))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "token "+gitInfo.PersonalAccessToken)
	req.Header.Set("Content-Type", "application/json")

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.StatusCode == 200 {
		fmt.Printf("Successfully applied label to PR# [%s].\n", prNumber)
	}
}
