package db

import (
	"time"
)

/*
	Date: 2023-10-04
	Description: Commit model
*/

type Author struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type Tree struct {
	Sha string `json:"sha"`
	URL string `json:"url"`
}

type Verification struct {
	Verified  bool   `json:"verified"`
	Reason    string `json:"reason"`
	Signature string `json:"signature"`
	Payload   string `json:"payload"`
}

type Commit struct {
	Author       Author       `json:"author"`
	Committer    Author       `json:"committer"`
	Message      string       `json:"message"`
	Tree         Tree         `json:"tree"`
	URL          string       `json:"url"`
	CommentCount int          `json:"comment_count"`
	Verification Verification `json:"verification"`
}

type Parent struct {
	Sha     string `json:"sha"`
	URL     string `json:"url"`
	HtmlURL string `json:"html_url"`
}

type GitHubCommit struct {
	Sha         string       `json:"sha"`
	NodeID      string       `json:"node_id"`
	Commit      Commit       `json:"commit"`
	URL         string       `json:"url"`
	HtmlURL     string       `json:"html_url"`
	CommentsURL string       `json:"comments_url"`
	Author      Contributors `json:"author"`
	Committer   Contributors `json:"committer"`
	Parents     []Parent     `json:"parents"`
}
