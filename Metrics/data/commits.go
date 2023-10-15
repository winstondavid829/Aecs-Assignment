package data

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type CommitsRequest struct {
	Owner          string `json:"Owner"`
	RepositoryName string `json:"RepositoryName"`
	// Author         string `json:"Author"`
	// Since          string `json:"Since"`
}

/*
	Date: 2023-10-01
	Description: Commit Response
*/

type User struct {
	AvatarURL         string `json:"avatar_url"`
	EventsURL         string `json:"events_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	GravatarID        string `json:"gravatar_id"`
	HTMLURL           string `json:"html_url"`
	ID                int    `json:"id"`
	Login             string `json:"login"`
	NodeID            string `json:"node_id"`
	OrganizationsURL  string `json:"organizations_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	ReposURL          string `json:"repos_url"`
	SiteAdmin         bool   `json:"site_admin"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	Type              string `json:"type"`
	URL               string `json:"url"`
}

type Tree struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type Verification struct {
	Payload   interface{} `json:"payload"`
	Reason    string      `json:"reason"`
	Signature interface{} `json:"signature"`
	Verified  bool        `json:"verified"`
}

type Commit struct {
	Author       Author       `json:"author"`
	CommentCount int          `json:"comment_count"`
	Committer    Author       `json:"committer"`
	Message      string       `json:"message"`
	Tree         Tree         `json:"tree"`
	URL          string       `json:"url"`
	Verification Verification `json:"verification"`
}

type Parent struct {
	HTMLURL string `json:"html_url"`
	SHA     string `json:"sha"`
	URL     string `json:"url"`
}

type GitHubCommit struct {
	Id          string    `json:"id"`
	Author      User      `json:"author"`
	CommentsURL string    `json:"comments_url"`
	Commit      Commit    `json:"commit"`
	Committer   User      `json:"committer"`
	HTMLURL     string    `json:"html_url"`
	NodeID      string    `json:"node_id"`
	Parents     []Parent  `json:"parents"`
	SHA         string    `json:"sha"`
	URL         string    `json:"url"`
	CreatedDate time.Time `json:"created_date"`
}

type GitHubCommitService interface {
	CreateCommits(m GitHubCommit) error
	InsertMany_Commits(commits []GitHubCommit) error
	CountCommitsByLogin(loginValue string) (int, error)
	// GetMovies() ([]Movie, error)
	// GetMovie(id string) (Movie, error)
	// UpdateMovie(m Movie) (Movie, error)
	// DeleteMovie(id string) error
}

func InitDatabase_Commits() GitHubCommitService {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &Database{
		client:    dynamodb.New(sess),
		tablename: "Github_Commits",
	}
}

func (db Database) CreateCommits(commits GitHubCommit) error {
	commits.Id = uuid.New().String()
	entityParsed, err := dynamodbattribute.MarshalMap(commits)
	if err != nil {
		return err
	}

	commits.CreatedDate = time.Now()

	input := &dynamodb.PutItemInput{
		Item:      entityParsed,
		TableName: aws.String(db.tablename),
	}

	_, err = db.client.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

/*
	Date: 2023-10-03
	Description: Insert Newly registered contributors
*/

func (db Database) InsertMany_Commits(commits []GitHubCommit) error {
	if len(commits) == 0 {
		return errors.New("No commits to insert")
	}

	for i := 0; i < len(commits); i += 25 {
		var writeRequests []*dynamodb.WriteRequest
		end := i + 25

		if end > len(commits) {
			end = len(commits)
		}

		for _, contributor := range commits[i:end] {
			// contributor.CreatedDate = time.Now()
			entityParsed, err := dynamodbattribute.MarshalMap(contributor)
			if err != nil {
				return err
			}

			writeRequests = append(writeRequests, &dynamodb.WriteRequest{
				PutRequest: &dynamodb.PutRequest{
					Item: entityParsed,
				},
			})
		}

		input := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				db.tablename: writeRequests,
			},
		}

		_, err := db.client.BatchWriteItem(input)
		if err != nil {
			return err
		}
	}

	return nil
}

/**/

func (db Database) CountCommitsByLoginAndCreatedAt(loginValue, createdAtDate string) (int, error) {
	// Convert map[string]types.AttributeValue to map[string]*dynamodb.AttributeValue
	attributeValues := make(map[string]*dynamodb.AttributeValue)
	attributeValues[":login"] = &dynamodb.AttributeValue{S: aws.String(loginValue)}
	// attributeValues[":startDate"] = &dynamodb.AttributeValue{S: aws.String(createdAtDate + "T00:00:00Z")}
	// attributeValues[":endDate"] = &dynamodb.AttributeValue{S: aws.String(createdAtDate + "T23:59:59Z")}

	// Define an expression attribute name for the reserved keyword "user.login"
	expressionAttributeNames := map[string]*string{
		"#user_login": aws.String("user.login"),
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(db.tablename),
		ExpressionAttributeValues: attributeValues,
		ExpressionAttributeNames:  expressionAttributeNames,
		FilterExpression: aws.String(
			"#user_login = :login"),
		// FilterExpression: aws.String(
		// 	"#user_login = :login AND created_at >= :startDate AND created_at <= :endDate"),
	}

	resp, err := db.client.Scan(input)
	if err != nil {
		fmt.Println("[CountCommitsByLoginAndCreatedAt]: Error querying DynamoDB:", err)
		return 0, err
	}

	count := len(resp.Items)
	fmt.Printf("Number of items with login %s and CreatedAt %s: %d\n", loginValue, createdAtDate, count)

	return count, nil
}

func (db Database) CountCommitsByLogin(loginValue string) (int, error) {
	// Convert map[string]types.AttributeValue to map[string]*dynamodb.AttributeValue
	attributeValues := make(map[string]*dynamodb.AttributeValue)
	attributeValues[":node_id"] = &dynamodb.AttributeValue{S: aws.String(loginValue)}

	// Define an expression attribute name for the reserved keyword "user.login"
	expressionAttributeNames := map[string]*string{
		"#user_login": aws.String("commit.node_id"),
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(db.tablename),
		ExpressionAttributeValues: attributeValues,
		ExpressionAttributeNames:  expressionAttributeNames,
		FilterExpression: aws.String(
			"#user_login = :node_id"),
	}

	resp, err := db.client.Scan(input)
	if err != nil {
		fmt.Println("[CountPullsByLogin]: Error scanning DynamoDB:", err)
		return 0, err
	}

	count := len(resp.Items)
	fmt.Printf("Number of items with login %s: %d\n", loginValue, count)

	return count, nil
}
