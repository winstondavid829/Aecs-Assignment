package data

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type UserMatrix struct {
	ID            string    `json:"id"`
	NodeID        string    `json:"NodeID"`
	Userlogin     string    `json:"Userlogin"`
	PullCount     int       `json:"PullCount"`
	CommitCount   int       `json:"CommitCount"`
	IssueComments int       `json:"IssueComments"`
	CreatedDate   time.Time `json:"CreatedDate"`
	ModifiedDate  time.Time `json:"ModifiedDate"`
}

type UserMatrixService interface {
	Create_UserMatrix(matrix UserMatrix) error
	GetMetricsByLoginID(userlogin string) ([]UserMatrix, error)
}

func InitDatabase_Matrix() UserMatrixService {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &Database{
		client:    dynamodb.New(sess),
		tablename: "UserMatrix",
	}
}

func (db Database) Create_UserMatrix(matrix UserMatrix) error {

	// contributor.CreatedDate = time.Now()
	matrix.ID = uuid.NewString()
	matrix.CreatedDate = time.Now()

	entityParsed, err := dynamodbattribute.MarshalMap(matrix)
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

func (db Database) GetMetricsByLoginID(userlogin string) ([]UserMatrix, error) {
	var metrics []UserMatrix
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue

	for {
		input := &dynamodb.ScanInput{
			TableName:                 aws.String(db.tablename),
			FilterExpression:          aws.String("Userlogin = :userlogin"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":userlogin": {S: aws.String(userlogin)}},
			ExclusiveStartKey:         lastEvaluatedKey,
		}

		result, err := db.client.Scan(input)
		if err != nil {
			return nil, err
		}

		for _, item := range result.Items {
			var metric UserMatrix
			if err := dynamodbattribute.UnmarshalMap(item, &metric); err != nil {
				return nil, err
			}
			metrics = append(metrics, metric)
		}

		lastEvaluatedKey = result.LastEvaluatedKey
		if len(lastEvaluatedKey) == 0 {
			break // No more results to fetch
		}
	}

	return metrics, nil
}
