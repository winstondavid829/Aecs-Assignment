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

type CommitsHandler struct {
	l *log.Logger
}

// NewUserHandler creates a new instance of UserHandler
func NewCommitsHandler(l *log.Logger) *CommitsHandler {
	return &CommitsHandler{l}
}

/*
Date: 2023-04-10
Description: Register Portpass Handler
*/
func (p *CommitsHandler) FetchGithub_CommitsHandlerV1(rw http.ResponseWriter, r *http.Request) {

	var requestbody data.CommitsRequest

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

	// Calculate the date for the previous day
	yesterday := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	today := time.Now().Format(time.RFC3339)

	// Update the API URL to include the date filter
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?&since=%s&until=%s",
		requestbody.Owner, requestbody.RepositoryName, yesterday, today)
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

	fmt.Printf("GitHub Commit Information:\n")
	fmt.Printf("OwnerName: %s\n", requestbody.Owner)
	fmt.Printf("RepositoryName: %s\n", requestbody.RepositoryName)
	// fmt.Printf("Author: %s\n", requestbody.Author)

	p.l.Println("Successfully retrieved commits details")
	rw.WriteHeader(http.StatusOK)
	responses := data.Responses{Status: true, Result: responseData}
	json.NewEncoder(rw).Encode(responses)
	return
}
