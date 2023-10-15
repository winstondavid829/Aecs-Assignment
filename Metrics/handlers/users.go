package handlers

import (
	"AECS_Assignment/thirdparty"
	"fmt"
	"os"
	"time"

	"github.com/sony/gobreaker"
)

var breaker *gobreaker.CircuitBreaker

func init() {
	// Create a Circuit Breaker with your desired settings
	breaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "GitHubAPI",
		MaxRequests: 3,               // Set an appropriate threshold
		Interval:    5 * time.Second, // Set an appropriate reset interval
	})
}

func Get_LoginUserData() {
	var response interface{}
	var err error

	// Use the Circuit Breaker for the API call
	response, err = breaker.Execute(func() (interface{}, error) {
		var apiIntegration thirdparty.APIRequestParams
		apiIntegration.ApiURL = "https://api.github.com/user"
		apiIntegration.AccessToken = os.Getenv("Github_Key")
		return apiIntegration.API_GetRequest()
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		// Handle the error, which may include fallback logic
		return
	}

	// Type assertion to map[string]interface{}
	responseData, ok := response.(map[string]interface{})
	if !ok {
		fmt.Printf("Unexpected response type: %T\n", response)
		// Handle the unexpected response type
		return
	}

	// Process the response
	fmt.Printf("GitHub User Information:\n")
	fmt.Printf("Username: %s\n", responseData["login"])
	fmt.Printf("Name: %s\n", responseData["name"])
	fmt.Printf("Location: %s\n", responseData["location"])
	fmt.Printf("Public Repositories: %.0f\n", responseData["public_repos"])
}

// var breaker *gobreaker.CircuitBreaker

// func init() {
// 	// Create a Circuit Breaker with your desired settings
// 	breaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
// 		Name:        "GitHubAPI",
// 		MaxRequests: 3,               // Set an appropriate threshold
// 		Interval:    5 * time.Second, // Set an appropriate reset interval
// 	})
// }

// func Get_Data() {
// 	var response interface{}
// 	var err error

// 	// Use the Circuit Breaker for the API call
// 	response, err = breaker.Execute(func() (interface{}, error) {
// 		var apiIntegration thirdparty.APIRequestParams
// 		apiIntegration.ApiURL = "https://api.github.com/user"
// 		apiIntegration.AccessToken = "https://api.github.com/user"
// 		return fetchDataFromGitHub()
// 	})

// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		// Handle the error, which may include fallback logic
// 		return
// 	}

// 	// Type assertion to map[string]interface{}
// 	responseData, ok := response.(map[string]interface{})
// 	if !ok {
// 		fmt.Printf("Unexpected response type: %T\n", response)
// 		// Handle the unexpected response type
// 		return
// 	}

// 	// Process the response
// 	fmt.Printf("GitHub User Information:\n")
// 	fmt.Printf("Username: %s\n", responseData["login"])
// 	fmt.Printf("Name: %s\n", responseData["name"])
// 	fmt.Printf("Location: %s\n", responseData["location"])
// 	fmt.Printf("Public Repositories: %.0f\n", responseData["public_repos"])
// }

// func fetchDataFromGitHub() (interface{}, error) {
// 	// Replace with your GitHub username and API token or personal access token
// 	// username := "your_username"
// 	token := os.Getenv("Github_Key")

// 	// Replace with the GitHub API URL you want to access
// 	apiURL := "https://api.github.com/user" // Example: Retrieve user information

// 	// Create an HTTP client
// 	client := &http.Client{}

// 	// Create an HTTP request
// 	req, err := http.NewRequest("GET", apiURL, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Add authorization header if using a personal access token
// 	if token != "" {
// 		req.Header.Add("Authorization", "token "+token)
// 	}

// 	// Send the HTTP request
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	// Check the response status code
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("Request failed with status: %s", resp.Status)
// 	}

// 	// Parse the JSON response
// 	var result map[string]interface{}
// 	err = json.NewDecoder(resp.Body).Decode(&result)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }
