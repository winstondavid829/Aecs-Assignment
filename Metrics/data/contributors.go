package data

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Database struct {
	client    *dynamodb.DynamoDB
	tablename string
}

type Contributors struct {
	ContributorId     int       `json:"id"`
	Login             string    `json:"login"`
	NodeID            string    `json:"node_id"`
	AvatarURL         string    `json:"avatar_url"`
	GRAvatarID        string    `json:"gravatar_id"`
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionURL   string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	RepoURL           string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Contributions     int       `json:"contributions"`
	CreatedDate       time.Time `json:"CreatedDate"`
	ModifiedDate      time.Time `json:"ModifiedDate"`
}

type ContributorsService interface {
	CreateContributors(m Contributors) error
	GetContributorbyID(id int) (Contributors, error)
	GetContributors() ([]Contributors, error)
	InsertMany_Contributors(contributors []Contributors) error
}

func InitDatabase(TableName string) ContributorsService {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &Database{
		client:    dynamodb.New(sess),
		tablename: TableName,
	}
}

func (db Database) CreateContributors(contributor Contributors) error {

	contributor.CreatedDate = time.Now()

	entityParsed, err := dynamodbattribute.MarshalMap(contributor)
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

func (db Database) GetContributors() ([]Contributors, error) {

	contributors := []Contributors{}

	// Initialize ScanInput without any filters or projections
	params := &dynamodb.ScanInput{
		TableName: aws.String(db.tablename),
	}

	// Perform the scan operation
	result, err := db.client.Scan(params)
	if err != nil {
		return []Contributors{}, err
	}

	// Iterate over the scan results and unmarshal into the Contributors struct
	for _, item := range result.Items {
		var contributor Contributors
		err = dynamodbattribute.UnmarshalMap(item, &contributor)
		if err != nil {
			return []Contributors{}, err
		}
		contributors = append(contributors, contributor)
	}

	return contributors, nil
}

func (db Database) GetContributorbyID(id int) (Contributors, error) {
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
		return Contributors{}, err
	}
	if result.Item == nil {
		msg := fmt.Sprintf("Contributor with id [ %v ] not found", id)
		return Contributors{}, errors.New(msg)
	}
	var contributor Contributors
	err = dynamodbattribute.UnmarshalMap(result.Item, &contributor)
	if err != nil {
		return Contributors{}, err
	}

	return contributor, nil
}

/*
	Date: 2023-10-03
	Description: Insert Newly registered contributors
*/

func (db Database) InsertMany_Contributors(contributors []Contributors) error {
	if len(contributors) == 0 {
		return errors.New("No contributors to insert")
	}

	for i := 0; i < len(contributors); i += 25 {
		var writeRequests []*dynamodb.WriteRequest
		end := i + 25

		if end > len(contributors) {
			end = len(contributors)
		}

		for _, contributor := range contributors[i:end] {
			contributor.CreatedDate = time.Now()
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
