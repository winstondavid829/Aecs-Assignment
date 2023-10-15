package data

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

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

type GithubPullService interface {
	CreatePulls(pulls GitHubPullRequest) error
	GetPulls() ([]GitHubPullRequest, error)
	GetPullsbyID(id int) (GitHubPullRequest, error)
	InsertMany_Pulls(pullrequests []GitHubPullRequest) error
	CountPullsByLogin(loginValue string) (int, error)
}

func InitDatabase_Pulls(TableName string) GithubPullService {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &Database{
		client:    dynamodb.New(sess),
		tablename: TableName,
	}
}

func (db Database) CreatePulls(pulls GitHubPullRequest) error {

	// contributor.CreatedDate = time.Now()

	entityParsed, err := dynamodbattribute.MarshalMap(pulls)
	if err != nil {
		return err
	}

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

func (db Database) GetPulls() ([]GitHubPullRequest, error) {

	pullRequests := []GitHubPullRequest{}

	// Initialize ScanInput without any filters or projections
	params := &dynamodb.ScanInput{
		TableName: aws.String(db.tablename),
	}

	// Perform the scan operation
	result, err := db.client.Scan(params)
	if err != nil {
		return []GitHubPullRequest{}, err
	}

	// Iterate over the scan results and unmarshal into the Contributors struct
	for _, item := range result.Items {
		var pull GitHubPullRequest
		err = dynamodbattribute.UnmarshalMap(item, &pull)
		if err != nil {
			return []GitHubPullRequest{}, err
		}
		pullRequests = append(pullRequests, pull)
	}

	return pullRequests, nil
}

func (db Database) GetPullsbyID(id int) (GitHubPullRequest, error) {
	fmt.Println("Contributor ID: ", id)

	result, err := db.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(db.tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				N: aws.String(strconv.Itoa(id)),
			},
		},
	})
	if err != nil {
		return GitHubPullRequest{}, err
	}
	if result.Item == nil {
		msg := fmt.Sprintf("Contributor with id [ %v ] not found", id)
		return GitHubPullRequest{}, errors.New(msg)
	}
	var pull GitHubPullRequest
	err = dynamodbattribute.UnmarshalMap(result.Item, &pull)
	if err != nil {
		return GitHubPullRequest{}, err
	}

	return pull, nil
}

/*
	Date: 2023-10-03
	Description: Insert Newly registered contributors
*/

func (db Database) InsertMany_Pulls(pullrequests []GitHubPullRequest) error {
	if len(pullrequests) == 0 {
		return errors.New("No contributors to insert")
	}

	for i := 0; i < len(pullrequests); i += 25 {
		var writeRequests []*dynamodb.WriteRequest
		end := i + 25

		if end > len(pullrequests) {
			end = len(pullrequests)
		}

		for _, contributor := range pullrequests[i:end] {
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

func (db Database) CountPullsByLogin(loginValue string) (int, error) {
	// Convert map[string]types.AttributeValue to map[string]*dynamodb.AttributeValue
	attributeValues := make(map[string]*dynamodb.AttributeValue)
	attributeValues[":login"] = &dynamodb.AttributeValue{S: aws.String(loginValue)}

	// Define an expression attribute name for the reserved keyword "user.login"
	expressionAttributeNames := map[string]*string{
		"#user_login": aws.String("author.login"),
	}

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(db.tablename),
		ExpressionAttributeValues: attributeValues,
		ExpressionAttributeNames:  expressionAttributeNames,
		FilterExpression: aws.String(
			"#user_login = :login"),
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
