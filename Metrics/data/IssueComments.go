package data

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Reactions struct {
	PlusOne    int    `json:"+1"`
	MinusOne   int    `json:"-1"`
	Confused   int    `json:"confused"`
	Eyes       int    `json:"eyes"`
	Heart      int    `json:"heart"`
	Hooray     int    `json:"hooray"`
	Laugh      int    `json:"laugh"`
	Rocket     int    `json:"rocket"`
	TotalCount int    `json:"total_count"`
	URL        string `json:"url"`
}

type IssueComment struct {
	AuthorAssociation     string      `json:"author_association"`
	Body                  string      `json:"body"`
	CreatedAt             string      `json:"created_at"`
	HTMLURL               string      `json:"html_url"`
	ID                    int         `json:"id"`
	IssueURL              string      `json:"issue_url"`
	NodeID                string      `json:"node_id"`
	PerformedViaGithubApp interface{} `json:"performed_via_github_app"` // null or some type
	Reactions             Reactions   `json:"reactions"`
	UpdatedAt             string      `json:"updated_at"`
	URL                   string      `json:"url"`
	User                  User        `json:"user"`
}

type GitHubCommitIssueComments interface {
	CreateIssues(issuecomments IssueComment) error
	InsertMany_IssueComments(comments []IssueComment) error
	CountIssuesByLogin(loginValue string) (int, error)
}

func InitDatabase_IssueComments() GitHubCommitIssueComments {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &Database{
		client:    dynamodb.New(sess),
		tablename: "Github_IssueComments",
	}
}

func (db Database) CreateIssues(issuecomments IssueComment) error {
	// commits.Id = uuid.New().String()
	entityParsed, err := dynamodbattribute.MarshalMap(issuecomments)
	if err != nil {
		fmt.Println("[CreateIssues]: ", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      entityParsed,
		TableName: aws.String(db.tablename),
	}

	_, err = db.client.PutItem(input)
	if err != nil {
		fmt.Println("[CreateIssues]: Failed inserting items", err)
		return err
	}

	return nil
}

/*
	Date: 2023-10-03
	Description: Insert Newly registered contributors
*/

func (db Database) InsertMany_IssueComments(comments []IssueComment) error {
	if len(comments) == 0 {
		return errors.New("No issue comments to insert")
	}

	for i := 0; i < len(comments); i += 25 {
		var writeRequests []*dynamodb.WriteRequest
		end := i + 25

		if end > len(comments) {
			end = len(comments)
		}

		for _, contributor := range comments[i:end] {
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

// func (db Database) CountIssuesByLoginAndCreatedAt(loginValue, createdAtDate string) (int, error) {
// 	// Convert map[string]types.AttributeValue to map[string]*dynamodb.AttributeValue
// 	attributeValues := make(map[string]*dynamodb.AttributeValue)
// 	attributeValues[":login"] = &dynamodb.AttributeValue{S: aws.String(loginValue)}
// 	attributeValues[":startDate"] = &dynamodb.AttributeValue{S: aws.String(createdAtDate + "T00:00:00Z")}
// 	attributeValues[":endDate"] = &dynamodb.AttributeValue{S: aws.String(createdAtDate + "T23:59:59Z")}

// 	// Define an expression attribute name for the reserved keyword "user"
// 	expressionAttributeNames := map[string]*string{
// 		"#user": aws.String("user"),
// 	}

// 	input := &dynamodb.ScanInput{
// 		TableName:                 aws.String(db.tablename),
// 		ExpressionAttributeValues: attributeValues,
// 		ExpressionAttributeNames:  expressionAttributeNames,
// 		FilterExpression: aws.String(
// 			"#user.login = :login AND created_at >= :startDate AND created_at <= :endDate"),
// 	}

// 	resp, err := db.client.Scan(input)
// 	if err != nil {
// 		fmt.Println("[CountIssuesByLoginAndCreatedAt]: Error querying DynamoDB:", err)
// 		return 0, err
// 	}

// 	count := len(resp.Items)
// 	fmt.Printf("Number of items with login %s and CreatedAt %s: %d\n", loginValue, createdAtDate, count)

// 	return count, nil
// }

func (db Database) CountIssuesByLogin(loginValue string) (int, error) {
	// Convert map[string]types.AttributeValue to map[string]*dynamodb.AttributeValue
	attributeValues := make(map[string]*dynamodb.AttributeValue)
	attributeValues[":login"] = &dynamodb.AttributeValue{S: aws.String(loginValue)}

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
	}

	resp, err := db.client.Scan(input)
	if err != nil {
		fmt.Println("[CountIssuesByLogin]: Error scanning DynamoDB:", err)
		return 0, err
	}

	count := len(resp.Items)
	fmt.Printf("Number of items with login %s: %d\n", loginValue, count)

	return count, nil
}
