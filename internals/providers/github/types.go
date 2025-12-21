package github

type User struct {
	GithubID int    `json:"github_id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}
