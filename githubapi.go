package gotip

type GitHubAPI []struct {
	Sha    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
	}
}
