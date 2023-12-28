package clients

import (
	"context"
	"github-events-microservices/collector/config"
	"github-events-microservices/collector/net"
	"github-events-microservices/model"
	"github.com/google/go-github/v57/github"
	"log/slog"
	"net/http"
)

type GitHubPublicEventsClient struct {
	restApiClient       GitHubRestClient
	githubGraphQLClient GithubGraphQLClient
}

func (receiver GitHubPublicEventsClient) ListEvents() ([]model.Event, error) {
	events, err := receiver.ListEventsWithOptions(&github.ListOptions{})
	return events, err
}

func (receiver GitHubPublicEventsClient) ListEventsWithOptions(options *github.ListOptions) ([]model.Event, error) {
	results, err := receiver.restApiClient.ListEvents(context.Background(), options)
	if err != nil {
		return nil, err
	}
	var events = make([]model.Event, len(results))
	for i, eventPointer := range results {
		events[i] = model.Event{
			ID:             *eventPointer.ID,
			Type:           *eventPointer.Type,
			CreatedAt:      eventPointer.CreatedAt.Time,
			Public:         *eventPointer.Public,
			RepoFullName:   *eventPointer.Repo.Name,
			RepoUrl:        *eventPointer.Repo.URL,
			ActorLogin:     *eventPointer.Actor.Login,
			ActorId:        *eventPointer.Actor.ID,
			ActorUrl:       *eventPointer.Actor.URL,
			ActorAvatarUrl: *eventPointer.Actor.AvatarURL,
		}
	}
	return events, nil
}

func (receiver GitHubPublicEventsClient) FetchRepos(events []model.Event) ([]model.Repo, error) {
	slog.Debug("fetching repos")
	return receiver.githubGraphQLClient.FetchRepos(events)
}

func NewGitHubClient(client *github.Client) *GitHubPublicEventsClient {
	graphQLClient := GithubGraphQLClient{
		httpClient:       net.NewSimpleHttpClient(http.Client{Timeout: config.CollectorConfiguration.GitHubGraphQLRequestTimeout}),
		githubGraphQLUrl: "https://api.github.com/graphql",
		token:            "Bearer " + config.CollectorConfiguration.GitHubToken,
	}

	return &GitHubPublicEventsClient{
		restApiClient:       NewSimpleGitHubRestClient(client),
		githubGraphQLClient: graphQLClient,
	}
}
