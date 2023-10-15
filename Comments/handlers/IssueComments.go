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

	"github.com/sony/gobreaker"
)

type IssueCommentsHandler struct {
	l *log.Logger
}

// NewUserHandler creates a new instance of UserHandler
func NewIssueCommentsHandler(l *log.Logger) *IssueCommentsHandler {
	return &IssueCommentsHandler{l}
}

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
Date: 2023-04-10
Description: Register Portpass Handler
*/
func (p *IssueCommentsHandler) FetchGithub_IssueCommentHandlerV1(rw http.ResponseWriter, r *http.Request) {

	var requestbody data.IssueCommentRequest

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

	// Calculate the date of the previous day
	yesterday := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	today := time.Now().Format(time.RFC3339)

	// Modify the API URL to include the 'since' and 'until' parameters
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/comments?since=%s&until=%s", requestbody.Owner, requestbody.RepositoryName, yesterday, today)

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

	// // Insert data into the database //

	fmt.Printf("GitHub Issue Comment Information:\n")
	fmt.Printf("OwnerName: %s\n", requestbody.Owner)
	fmt.Printf("RepositoryName: %s\n", requestbody.RepositoryName)

	// _, err1 := mongodb.CreateContainers(requestbody)

	// if err1 != nil {
	// 	p.l.Println("Failed adding containers", err1)
	// 	rw.WriteHeader(http.StatusInternalServerError)
	// 	responses := data.Responses{Status: false, Result: fmt.Sprintf("Failed adding containers: %v", err1)}
	// 	json.NewEncoder(rw).Encode(responses)
	// 	return
	// }

	p.l.Println("Successfully fetched Github Issue Comment Information")
	rw.WriteHeader(http.StatusOK)
	responses := data.Responses{Status: true, Result: responseData}
	json.NewEncoder(rw).Encode(responses)
	return
}
