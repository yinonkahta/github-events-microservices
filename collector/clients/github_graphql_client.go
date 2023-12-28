package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github-events-microservices/collector/net"
	"github-events-microservices/model"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	repoQueryPrefix = "repoQuery"
)

type GithubGraphQLClient struct {
	httpClient       net.HttpClient
	githubGraphQLUrl string
	token            string
}

type RepoIdentifier struct {
	Owner string
	Name  string
}

type GraphQLQuery struct {
	Query string `json:"query"`
}

type RepoQueryResponse struct {
	Data   map[string]RepoData    `json:"data"`
	Errors map[string]interface{} `json:"errors"`
}

type RepoData struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Owner RepoOwner `json:"owner"`
	Url   string    `json:"url"`
	Stars int       `json:"stargazerCount"`
}

type RepoOwner struct {
	Login string `json:"login"`
}

func (receiver GithubGraphQLClient) FetchRepos(events []model.Event) ([]model.Repo, error) {
	repoLastUpdatedMap := buildRepoLastUpdatedMap(events)
	query := buildQuery(events)
	body, err := receiver.sendRequest(query)
	if err != nil {
		return nil, err
	}

	repos, err := parseQueryResponse(body, repoLastUpdatedMap)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func (receiver GithubGraphQLClient) sendRequest(query GraphQLQuery) ([]byte, error) {
	marshal, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", receiver.githubGraphQLUrl, bytes.NewReader(marshal))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create fetch repos request: %s", err.Error()))
		return nil, err
	}

	request.Header.Add("Authorization", receiver.token)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to fetch repos: %s", err.Error()))
	}

	response, err := receiver.httpClient.Do(request)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to fetch repos: %s", err.Error()))
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	slog.Debug(fmt.Sprintf("response Body: %s", string(body)))
	return body, nil
}

func buildRepoLastUpdatedMap(events []model.Event) map[RepoIdentifier]time.Time {
	repoToLastUpdated := make(map[RepoIdentifier]time.Time)
	for _, event := range events {
		repoIdentifier, err := getRepoIdentifier(event.RepoFullName)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to parse repo identifier from: %s. Reason: %s", event.RepoFullName, err.Error()))
			continue
		}
		lastUpdated, ok := repoToLastUpdated[*repoIdentifier]
		if !ok || (ok && event.CreatedAt.After(lastUpdated)) {
			repoToLastUpdated[*repoIdentifier] = event.CreatedAt
		}
	}
	return repoToLastUpdated
}

func getRepoIdentifier(repoName string) (*RepoIdentifier, error) {
	repoOwnerAndName := strings.Split(repoName, "/")
	if len(repoOwnerAndName) != 2 || len(repoOwnerAndName[0]) == 0 || len(repoOwnerAndName[0]) == 0 {
		return nil, errors.New(fmt.Sprintf("Failed to parse repo identifier from: %s", repoName))
	}

	return &RepoIdentifier{
		Owner: repoOwnerAndName[0],
		Name:  repoOwnerAndName[1],
	}, nil
}

func buildQuery(events []model.Event) GraphQLQuery {
	var sb strings.Builder
	sb.WriteString("query {")
	for i, event := range events {
		repoOwnerAndName := strings.Split(event.RepoFullName, "/")
		sb.WriteString(fmt.Sprintf("%s%d:repository(owner: \"%s\", name: \"%s\") {id,name,owner{login},url,stargazerCount}", repoQueryPrefix, i, repoOwnerAndName[0], repoOwnerAndName[1]))
	}
	sb.WriteString("}")

	queryString := sb.String()
	return GraphQLQuery{Query: queryString}
}

func parseQueryResponse(bytes []byte, repoToLastUpdated map[RepoIdentifier]time.Time) ([]model.Repo, error) {
	var repoQueryResponse RepoQueryResponse
	err := json.Unmarshal(bytes, &repoQueryResponse)
	if err != nil {
		return nil, err
	}

	if repoQueryResponse.Errors != nil {
		logErrors(repoQueryResponse.Errors)
	}

	var repos = make([]model.Repo, 0)
	for _, repoData := range repoQueryResponse.Data {
		repoIdentifier := RepoIdentifier{
			Owner: repoData.Owner.Login,
			Name:  repoData.Name,
		}
		repo := model.Repo{
			ID:            repoData.ID,
			Owner:         repoIdentifier.Owner,
			Name:          repoIdentifier.Name,
			Url:           repoData.Url,
			Stars:         repoData.Stars,
			LastUpdatedAt: repoToLastUpdated[repoIdentifier],
		}
		repos = append(repos, repo)
	}
	return repos, nil
}

func logErrors(responseMap map[string]interface{}) {
	errorsData, ok := responseMap["errors"]
	if !ok {
		return
	}
	errorsSlice := errorsData.([]interface{})
	errorMessages := make([]string, 0)
	for _, errorVal := range errorsSlice {
		errorMessage, ok := errorVal.(map[string]interface{})["message"]
		if ok {
			errorMessages = append(errorMessages, errorMessage.(string))
		}
	}

	if len(errorMessages) > 0 {
		slog.Error(fmt.Sprintf("Failed to fetch the following repos due to the following errors: %s", strings.Join(errorMessages, ", ")))
	}
}
