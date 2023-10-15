package data

type CommitsRequest struct {
	Owner          string `json:"Owner"`
	RepositoryName string `json:"RepositoryName"`
	// Author         string `json:"Author"`
	// Since          string `json:"Since"`
}
