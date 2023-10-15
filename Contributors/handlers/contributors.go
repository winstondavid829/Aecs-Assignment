package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Collaborators/constants"
	"github.com/Collaborators/data"
	db "github.com/Collaborators/dbdocs"
	"github.com/Collaborators/thirdparty"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sony/gobreaker"
)

type ContributorsHandler struct {
	l *log.Logger
}

// NewUserHandler creates a new instance of UserHandler
func NewContributorsHandler(l *log.Logger) *ContributorsHandler {
	return &ContributorsHandler{l}
}

const Owner = "MrB141107"
const RepositoryName = "Hacktoberfest_2022"

/*
Date: 2023-04-10
Description: Register Portpass Handler
*/

var breaker *gobreaker.CircuitBreaker

func init() {
	// Create a Circuit Breaker with your desired settings
	breaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "GitHubAPI",
		MaxRequests: 3,               // Set an appropriate threshold
		Interval:    5 * time.Second, // Set an appropriate reset interval
	})
}

/*
	Date: 2023-10-04
	Description: Upload data to S3 bucket
*/

func uploadToS3(data []byte, fileName string) {
	// Initialize a session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-1"),
	}))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("aecsdata"),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		fmt.Printf("Failed to upload file, %v", err)
		return
	}
	fmt.Printf("Successfully uploaded %s to S3\n", fileName)
}

/*
	Date: 2023-10-04
	Description: Upload data to S3 bucket
*/

func filterByPreviousDay(data []map[string]interface{}) []map[string]interface{} {
	filteredData := []map[string]interface{}{}
	prevDay := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	for _, record := range data {
		if createdAt, ok := record["created_at"].(string); ok {
			if strings.Contains(createdAt, prevDay) {
				filteredData = append(filteredData, record)
			}
		}
	}
	return filteredData
}

/*
	Date: 2023-10-04
	Description: Function to fetch data from API
*/

func fetchDataFromAPI(apiUrl string) []map[string]interface{} {
	// Use the Circuit Breaker for the API call
	response, err := breaker.Execute(func() (interface{}, error) {
		var apiIntegration thirdparty.APIRequestParams
		apiIntegration.ApiURL = apiUrl
		apiIntegration.AccessToken = os.Getenv("Github_Key")
		return apiIntegration.API_GetRequest()
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		// Handle the error, which may include fallback logic
		return nil
	}

	// Type assertion to map[string]interface{}
	responseData, ok := response.([]map[string]interface{})
	if !ok {
		fmt.Printf("Unexpected response type: %T\n", response)
		// Handle the unexpected response type
		return nil
	}

	return responseData
}

/*
	Date: 2023-10-04
	Description: Function to fetch Repo Contributor List
*/

func (p *ContributorsHandler) FetchGithubRepo_Contributors() {
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors", Owner, RepositoryName)
	// Fetch and filter data
	rawData := fetchDataFromAPI(apiUrl) // Assume fetchDataFromAPI is a function that fetches data from GitHub API
	filteredData := filterByPreviousDay(rawData)
	// Upload to S3
	// jsonData, _ := json.Marshal(filteredData)
	// uploadToS3(jsonData, "contributors_"+time.Now().Format("2006-01-02")+".json")

	// Marshal slice of maps to JSON
	jsonData, err := json.Marshal(filteredData)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	// Initialize slice of structs and Unmarshal JSON to it
	var contributors []db.Contributors
	err01 := json.Unmarshal(jsonData, &contributors)
	if err01 != nil {
		p.l.Println("Error unmarshalling to struct:", err01)
		return
	}

	err02 := db.InitDatabase().InsertMany_Contributors(contributors)

	if err02 != nil {
		p.l.Println("Error inserting contributors record", err02)
		return
	}
}

/*
	Date: 2023-10-04
	Description: Function to fetch Repo Pull Request List
*/

func (p *ContributorsHandler) FetchGithubRepo_Pulls() {
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", Owner, RepositoryName)
	// Fetch and filter data
	rawData := fetchDataFromAPI(apiUrl)
	filteredData := filterByPreviousDay(rawData)
	// Upload to S3
	jsonData, _ := json.Marshal(filteredData)
	uploadToS3(jsonData, "pulls_"+time.Now().Format("2006-01-02")+".json")
}

/*
	Date: 2023-10-04
	Description: Function to fetch Repo Commits List
*/

func (p *ContributorsHandler) FetchGithubRepo_Commits() {
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", Owner, RepositoryName)
	// Fetch and filter data
	rawData := fetchDataFromAPI(apiUrl)
	filteredData := filterByPreviousDay(rawData)
	// Upload to S3
	jsonData, _ := json.Marshal(filteredData)
	uploadToS3(jsonData, "commits_"+time.Now().Format("2006-01-02")+".json")
}

/*
	Date: 2023-10-04
	Description: Function to fetch Repo Issues Comment List
*/

func (p *ContributorsHandler) FetchGithubRepo_Comments() {
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/comments", Owner, RepositoryName)
	// Fetch and filter data
	rawData := fetchDataFromAPI(apiUrl)
	filteredData := filterByPreviousDay(rawData)
	// Upload to S3
	jsonData, _ := json.Marshal(filteredData)
	uploadToS3(jsonData, "comments_"+time.Now().Format("2006-01-02")+".json")
}

/*
	Date: 2023-10-09
	Description: Fetch yesterday joined new contributors
*/

func (p *ContributorsHandler) FetchGithub_ContributorsHandlerV1(rw http.ResponseWriter, r *http.Request) {

	var requestbody data.CommitsRequest

	// validate header //
	// if auth.ValidateUserTokenInHeader(r) {
	// 	p.l.Println("Failed validating header token")
	// 	rw.WriteHeader(http.StatusBadRequest)
	// 	responses := data.Responses{Status: false, Result: fmt.Sprintf("%v", "Failed validating header token")}
	// 	json.NewEncoder(rw).Encode(responses)
	// 	return
	// }

	// Deserialize the request body //
	err := json.NewDecoder(r.Body).Decode(&requestbody)
	if err != nil {
		p.l.Println("Deserialization error:", err)
		rw.WriteHeader(http.StatusBadRequest)
		responses := data.Responses{Status: false, Result: fmt.Sprintf("Deserialization error: %v", err)}
		json.NewEncoder(rw).Encode(responses)
		return
	}

	// Validate user input
	if err := constants.Validate.Struct(requestbody); err != nil {
		p.l.Println("Validation error:", err)
		rw.WriteHeader(http.StatusBadRequest)
		responses := data.Responses{Status: false, Result: fmt.Sprintf("Deserialization error: %v", err)}
		json.NewEncoder(rw).Encode(responses)
		return
	}

	// Construct the API URL with since and until query parameters
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors", requestbody.Owner, requestbody.RepositoryName)

	// Use the Circuit Breaker for the API call
	response, err := breaker.Execute(func() (interface{}, error) {
		var apiIntegration thirdparty.APIRequestParams
		apiIntegration.ApiURL = apiUrl
		apiIntegration.AccessToken = os.Getenv("Github_Key")
		return apiIntegration.API_GetRequest()
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		// Handle the error, which may include fallback logic
		return
	}

	// Type assertion to map[string]interface{}
	responseData, ok := response.([]map[string]interface{})
	if !ok {
		fmt.Printf("Unexpected response type: %T\n", response)
		// Handle the unexpected response type
		return
	}

	// Filter data before inserting //
	// filteredData := filterByPreviousDay(responseData)

	fmt.Printf("GitHub Commit Information:\n")
	fmt.Printf("OwnerName: %s\n", requestbody.Owner)
	fmt.Printf("RepositoryName: %s\n", requestbody.RepositoryName)

	p.l.Println("Successfully added new container details")
	rw.WriteHeader(http.StatusOK)
	responses := data.Responses{Status: true, Result: responseData}
	json.NewEncoder(rw).Encode(responses)
	return
}
