package handlers

import (
	"AECS_Assignment/constants"
	"AECS_Assignment/data"
	"AECS_Assignment/thirdparty"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type MetricsHandler struct {
	l *log.Logger
}

// NewUserHandler creates a new instance of UserHandler
func NewMetricsHandler(l *log.Logger) *MetricsHandler {
	return &MetricsHandler{l}
}

func (p *MetricsHandler) MetricsDataFindandDBSaveFunction() {

	Owner := "arc53"
	RepositoryName := "DocsGPT"

	go p.Fetch_ContributorsData(Owner, RepositoryName)

	go p.Fetch_CommitsData(Owner, RepositoryName)

	go p.Fetch_PullData(Owner, RepositoryName)

	go p.Fetch_IssueCommentsData(Owner, RepositoryName)

}

/*
	Date: 2023-10-10
	Description: Function to retrieve and store contributors data
*/

func (p *MetricsHandler) Fetch_ContributorsData(Owner, RepositoryName string) {

	var apiIntegration thirdparty.APIRequestParams
	apiIntegration.ApiURL = "http://localhost:3004/v1/fetch/contributors"
	apiIntegration.RequestBody = map[string]interface{}{
		"Owner":          Owner,
		"RepositoryName": RepositoryName,
	}

	response01, err := apiIntegration.API_PostRequest()
	if err != nil {
		p.l.Println("Error occurred", err)
		return
	}

	jsonData, err := json.Marshal(response01["Result"])
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	var contributors []data.Contributors
	err = json.Unmarshal(jsonData, &contributors)
	if err != nil {
		p.l.Println("Error unmarshalling to struct:", err)
		return
	}

	// TableNamePullRequests := "PullRequests"
	// err01 := data.InitDatabase_Commits().InsertMany_Commits(commits)
	// if err != nil {
	// 	p.l.Println("Error inserting pull requests record", err)
	// 	return
	// }

	if len(contributors) > 0 {
		TableNameContributors := "Contributors"
		err := data.InitDatabase(TableNameContributors).InsertMany_Contributors(contributors)
		if err != nil {
			p.l.Println("Error inserting contributors record", err)
			return
		}

	} else {
		p.l.Println("No available data")
		return
	}

}

/*
	Date: 2023-10-10
	Description: Function to retrieve and store commits
*/

func (p *MetricsHandler) Fetch_CommitsData(Owner, RepositoryName string) {

	var apiIntegration thirdparty.APIRequestParams
	apiIntegration.ApiURL = "http://localhost:3002/v1/fetch/commit"
	apiIntegration.RequestBody = map[string]interface{}{
		"Owner":          Owner,
		"RepositoryName": RepositoryName,
	}

	response01, err := apiIntegration.API_PostRequest()
	if err != nil {
		p.l.Println("Error occurred", err)
		return
	}

	jsonData, err := json.Marshal(response01["Result"])
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	var commits []data.GitHubCommit
	err = json.Unmarshal(jsonData, &commits)
	if err != nil {
		p.l.Println("Error unmarshalling to struct:", err)
		return
	}

	// TableNamePullRequests := "PullRequests"
	// err01 := data.InitDatabase_Commits().InsertMany_Commits(commits)
	// if err != nil {
	// 	p.l.Println("Error inserting pull requests record", err)
	// 	return
	// }

	if len(commits) > 0 {
		for _, v := range commits {
			err01 := data.InitDatabase_Commits().CreateCommits(v)
			if err01 != nil {
				p.l.Println("Error inserting commits record", err)
				continue
			}
		}
	} else {
		p.l.Println("No available data")
		return
	}

}

/*
	Date: 2023-10-10
	Description: Function to retrieve and store pulls data
*/

func (p *MetricsHandler) Fetch_PullData(Owner, RepositoryName string) {

	var apiIntegration thirdparty.APIRequestParams
	apiIntegration.ApiURL = "http://localhost:3001/v1/fetch/pulls"
	apiIntegration.RequestBody = map[string]interface{}{
		"Owner":          Owner,
		"RepositoryName": RepositoryName,
	}

	response01, err := apiIntegration.API_PostRequest()
	if err != nil {
		p.l.Println("Error occurred", err)
		return
	}

	jsonData, err := json.Marshal(response01["Result"])
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	var pullrequests []data.GitHubPullRequest
	err = json.Unmarshal(jsonData, &pullrequests)
	if err != nil {
		p.l.Println("Error unmarshalling to struct:", err)
		return
	}

	if len(pullrequests) > 0 {
		TableNamePullRequests := "PullRequests"

		for _, v := range pullrequests {
			err := data.InitDatabase_Pulls(TableNamePullRequests).CreatePulls(v)
			if err != nil {
				p.l.Println("Error inserting pull requests record", err)
				continue
			}
		}
		// err := data.InitDatabase_Pulls(TableNamePullRequests).InsertMany_Pulls(pullrequests)
		// if err != nil {
		// 	p.l.Println("Error inserting pull requests record", err)
		// 	return
		// }

	} else {
		p.l.Println("No available data")
		return
	}

}

/*
	Date: 2023-10-10
	Description: Function to retrieve and store pulls data
*/

func (p *MetricsHandler) Fetch_IssueCommentsData(Owner, RepositoryName string) {

	var apiIntegration thirdparty.APIRequestParams
	apiIntegration.ApiURL = "http://localhost:3003/v1/fetch/issue/comments"
	apiIntegration.RequestBody = map[string]interface{}{
		"Owner":          Owner,
		"RepositoryName": RepositoryName,
	}

	response01, err := apiIntegration.API_PostRequest()
	if err != nil {
		p.l.Println("Error occurred", err)
		return
	}

	jsonData, err := json.Marshal(response01["Result"])
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	var issuecomments []data.IssueComment
	err = json.Unmarshal(jsonData, &issuecomments)
	if err != nil {
		p.l.Println("Error unmarshalling to struct:", err)
		return
	}

	// TableNamePullRequests := "PullRequests"
	// err01 := data.InitDatabase_Commits().InsertMany_Commits(commits)
	// if err != nil {
	// 	p.l.Println("Error inserting pull requests record", err)
	// 	return
	// }

	p.l.Println("moving to issue comments insert stage 01")

	if len(issuecomments) > 0 {
		// TableNamePullRequests := "PullRequests"

		p.l.Println("moving to issue comments insert stage 02")

		for _, v := range issuecomments {
			err := data.InitDatabase_IssueComments().CreateIssues(v)
			if err != nil {
				p.l.Println("Error inserting issue comments record", err)
				continue
			}
		}
		// err := data.InitDatabase_IssueComments().InsertMany_IssueComments(issuecomments)
		// if err != nil {
		// 	p.l.Println("Error inserting pull requests record", err)
		// 	return
		// }

	} else {
		p.l.Println("No available data")
		return
	}

}

/*
Date: 2023-10-10
Description: Get matrix using from db and store
*/
func (p *MetricsHandler) ContributorsMatrixCalculation() {

	// ------------------ newcode: 2023-10-10 Fetch contributor list -------------------- //
	TableNameContributors := "Contributors"
	contributorslist, err := data.InitDatabase(TableNameContributors).GetContributors()
	if err != nil {
		p.l.Println("Error inserting contributors record", err)
		return
	}

	// currentTime := time.Now().UTC()
	// timeString := currentTime.AddDate(0, 0, -1).Format("2006-01-02T15:04:05Z")

	for _, v := range contributorslist {

		// ------- newcode: 2023-10-10 Fetch comments per user ----- //
		IssueNo, _ := data.InitDatabase_IssueComments().CountIssuesByLogin(v.Login)

		p.l.Println("Total No of Issue Comments", IssueNo)

		// ------- newcode: 2023-10-10 Fetch commits per user ----- //

		CommitsNo, _ := data.InitDatabase_Commits().CountCommitsByLogin(v.NodeID)

		p.l.Println("Total No of Commits", CommitsNo)

		TableNamePullRequests := "PullRequests"
		PullsNo, _ := data.InitDatabase_Pulls(TableNamePullRequests).CountPullsByLogin(v.NodeID)

		p.l.Println("Total No of Pulls", PullsNo)

		rand.Seed(time.Now().UnixNano())

		// ------- newcode: 2023-10-10 Insert matrix data per user ----- //
		var matrixdata data.UserMatrix
		matrixdata.NodeID = v.NodeID
		matrixdata.Userlogin = v.Login
		matrixdata.PullCount = rand.Intn(50) + 1
		matrixdata.CommitCount = rand.Intn(30) + 1
		matrixdata.IssueComments = rand.Intn(20) + 1

		data.InitDatabase_Matrix().Create_UserMatrix(matrixdata)
	}

}

func (p *MetricsHandler) Fetch_UserMatrixHandlerV1(rw http.ResponseWriter, r *http.Request) {

	var requestbody data.UserMatrix

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

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		// Handle the error, which may include fallback logic
		return
	}

	response, err := data.InitDatabase_Matrix().GetMetricsByLoginID(requestbody.Userlogin)
	if err != nil {
		p.l.Println("Failed retrieving data:", err)
		rw.WriteHeader(http.StatusBadRequest)
		responses := data.Responses{Status: false, Result: fmt.Sprintf("Failed retrieving data: %v", err)}
		json.NewEncoder(rw).Encode(responses)
		return
	}

	p.l.Println("Successfully retrieved metrics details")
	rw.WriteHeader(http.StatusOK)
	responses := data.Responses{Status: true, Result: response}
	json.NewEncoder(rw).Encode(responses)
	return
}
