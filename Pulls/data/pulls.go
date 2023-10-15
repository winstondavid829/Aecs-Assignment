package data

type CommitsRequest struct {
	Owner          string `json:"Owner"`
	RepositoryName string `json:"RepositoryName"`
}

/*
	Date:2023-10-04
	Description: Pull Requests model
*/

type Author struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

type RepoDetails struct {
	FullName string `json:"full_name"`
	Language string `json:"language"`
}

type URLs struct {
	HTMLURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
}

type GitHubPullRequest struct {
	PRID        int64       `json:"pr_id"`
	Title       string      `json:"title"`
	State       string      `json:"state"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Author      Author      `json:"author"`
	Labels      interface{} `json:"labels"`
	RepoDetails RepoDetails `json:"repo_details"`
	URLs        URLs        `json:"urls"`
	Body        string      `json:"body"`
}
