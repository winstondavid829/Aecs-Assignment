package handlers

import (
	"AECS_Assignment/constants"
	"AECS_Assignment/data"
	"AECS_Assignment/thirdparty"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type PullsHandler struct {
	l *log.Logger
}

// NewUserHandler creates a new instance of UserHandler
func NewPullsHandler(l *log.Logger) *PullsHandler {
	return &PullsHandler{l}
}

/*
	Date: 2023-10-09
	Description: Pull Request using cron
*/

func (p *PullsHandler) FetchGithub_PullsV1() {

	Owner := "arc53"
	RepositoryName := "DocsGPT"

	// Get the current date and time
	currentTime := time.Now()

	// Get the previous day's date
	previousDay := currentTime.AddDate(0, 0, -1).Format("2006-01-02T15:04:05Z")

	// Get the current day's date
	currentDay := currentTime.Format("2006-01-02T15:04:05Z")

	// Construct the API URL with since and until query parameters
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?since=%s&until=%s", Owner, RepositoryName, previousDay, currentDay)

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

	// Marshal slice of maps to JSON
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	// Initialize slice of structs and Unmarshal JSON to it
	var pullRequests []data.GitHubPullRequest
	err01 := json.Unmarshal(jsonData, &pullRequests)
	if err01 != nil {
		p.l.Println("Error unmarshalling to struct:", err01)
		return
	}

	b, err := json.Marshal(pullRequests)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	fmt.Printf("GitHub Pull Request Information:\n")
	fmt.Printf("OwnerName: %s\n", Owner)
	fmt.Printf("RepositoryName: %s\n", RepositoryName)
	fmt.Printf("Pulls from previous day: %v\n", responseData)
}

/*
Date: 2023-04-10
Description: Register Portpass Handler
*/
func (p *PullsHandler) FetchGithub_PullsHandlerV1(rw http.ResponseWriter, r *http.Request) {

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

	// Get the current date and time
	currentTime := time.Now()

	// Get the previous day's date
	previousDay := currentTime.AddDate(0, 0, -1).Format("2006-01-02T15:04:05Z")

	// Get the current day's date
	currentDay := currentTime.Format("2006-01-02T15:04:05Z")

	// Construct the API URL with since and until query parameters
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?since=%s&until=%s", requestbody.Owner, requestbody.RepositoryName, previousDay, currentDay)

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

	// Marshal slice of maps to JSON
	jsonData, err := json.Marshal(responseData)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	// Initialize slice of structs and Unmarshal JSON to it
	var pullRequests []data.GitHubPullRequest
	err01 := json.Unmarshal(jsonData, &pullRequests)
	if err01 != nil {
		p.l.Println("Error unmarshalling to struct:", err01)
		return
	}

	// b, err := json.Marshal(pullRequests)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(b))

	// // Insert data into the database //

	fmt.Printf("GitHub Commit Information:\n")
	fmt.Printf("OwnerName: %s\n", requestbody.Owner)
	fmt.Printf("RepositoryName: %s\n", requestbody.RepositoryName)

	p.l.Println("Successfully added new container details")
	rw.WriteHeader(http.StatusOK)
	responses := data.Responses{Status: true, Result: responseData}
	json.NewEncoder(rw).Encode(responses)
	return
}
