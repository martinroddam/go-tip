package gotip

type GitHubAPI []GitHubCommit

type GitHubCommit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
	}
}
