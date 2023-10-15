package data

type IssueCommentRequest struct {
	Owner          string `json:"Owner"`
	RepositoryName string `json:"RepositoryName"`
}
