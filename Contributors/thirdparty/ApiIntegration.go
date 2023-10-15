package thirdparty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type APIRequestParams struct {
	AccessToken string
	ApiURL      string
	RequestBody map[string]interface{}
}

func (p *APIRequestParams) API_GetRequest() (interface{}, error) {
	// Replace with your GitHub username and API token or personal access token
	// username := "your_username"
	token := p.AccessToken

	// Replace with the GitHub API URL you want to access
	apiURL := p.ApiURL // Example: Retrieve user information

	// Create an HTTP client
	client := &http.Client{}

	// Create an HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Add authorization header if using a personal access token
	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request failed with status: %s", resp.Status)
	}

	// Parse the JSON response
	var result []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *APIRequestParams) API_PostRequest() (interface{}, error) {
	token := p.AccessToken
	apiURL := p.ApiURL

	client := &http.Client{}

	// Convert the request body to JSON
	jsonBody, err := json.Marshal(p.RequestBody)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf("Request failed with status: %s, Body: %s", resp.Status, bodyString)
	}

	// Parse the JSON response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
