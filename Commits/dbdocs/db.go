package db

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type Database struct {
	client    *dynamodb.DynamoDB
	tablename string
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

type Author struct {
	Date  string `json:"date"`
	Email string `json:"email"`
	Name  string `json:"name"`
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
	Id          string   `json:"id"`
	Author      User     `json:"author"`
	CommentsURL string   `json:"comments_url"`
	Commit      Commit   `json:"commit"`
	Committer   User     `json:"committer"`
	HTMLURL     string   `json:"html_url"`
	NodeID      string   `json:"node_id"`
	Parents     []Parent `json:"parents"`
	SHA         string   `json:"sha"`
	URL         string   `json:"url"`
}

type GitHubCommitService interface {
	CreateCommits(m GitHubCommit) error
	// GetMovies() ([]Movie, error)
	// GetMovie(id string) (Movie, error)
	// UpdateMovie(m Movie) (Movie, error)
	// DeleteMovie(id string) error
}

func InitDatabase() GitHubCommitService {
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

// func (db Database) GetMovies() ([]Movie, error) {
// 	movies := []Movie{}
// 	filt := expression.Name("Id").AttributeNotExists()
// 	proj := expression.NamesList(
// 		expression.Name("id"),
// 		expression.Name("name"),
// 	)
// 	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
// 	if err != nil {
// 		return []Movie{}, err
// 	}
// 	params := &dynamodb.ScanInput{
// 		ExpressionAttributeNames:  expr.Names(),
// 		ExpressionAttributeValues: expr.Values(),
// 		FilterExpression:          expr.Filter(),
// 		ProjectionExpression:      expr.Projection(),
// 		TableName:                 aws.String(db.tablename),
// 	}
// 	result, err := db.client.Scan(params)

// 	if err != nil {

// 		return []Movie{}, err
// 	}

// 	for _, item := range result.Items {
// 		var movie Movie
// 		err = dynamodbattribute.UnmarshalMap(item, &movie)
// 		movies = append(movies, movie)

// 	}

// 	return movies, nil
// }

// func (db Database) GetMovie(id string) (Movie, error) {
// 	result, err := db.client.GetItem(&dynamodb.GetItemInput{
// 		TableName: aws.String(db.tablename),
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"id": {
// 				S: aws.String(id),
// 			},
// 		},
// 	})
// 	if err != nil {
// 		return Movie{}, err
// 	}
// 	if result.Item == nil {
// 		msg := fmt.Sprintf("Movie with id [ %s ] not found", id)
// 		return Movie{}, errors.New(msg)
// 	}
// 	var movie Movie
// 	err = dynamodbattribute.UnmarshalMap(result.Item, &movie)
// 	if err != nil {
// 		return Movie{}, err
// 	}

// 	return movie, nil
// }

// func (db Database) UpdateMovie(movie Movie) (Movie, error) {
// 	entityParsed, err := dynamodbattribute.MarshalMap(movie)
// 	if err != nil {
// 		return Movie{}, err
// 	}

// 	input := &dynamodb.PutItemInput{
// 		Item:      entityParsed,
// 		TableName: aws.String(db.tablename),
// 	}

// 	_, err = db.client.PutItem(input)
// 	if err != nil {
// 		return Movie{}, err
// 	}

// 	return movie, nil
// }

// func (db Database) DeleteMovie(id string) error {
// 	input := &dynamodb.DeleteItemInput{
// 		Key: map[string]*dynamodb.AttributeValue{
// 			"id": {
// 				S: aws.String(id),
// 			},
// 		},
// 		TableName: aws.String(db.tablename),
// 	}

// 	res, err := db.client.DeleteItem(input)
// 	if res == nil {
// 		return errors.New(fmt.Sprintf("No movie to de: %s", err))
// 	}
// 	if err != nil {
// 		return errors.New(fmt.Sprintf("Got error calling DeleteItem: %s", err))
// 	}
// 	return nil
// }
