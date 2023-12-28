package clients

import (
	"context"
	"github.com/google/go-github/v57/github"
)

type GitHubRestClient interface {
	ListEvents(context.Context, *github.ListOptions) ([]*github.Event, error)
}

type SimpleGitHubRestClient struct {
	restApiClient *github.Client
}

func NewSimpleGitHubRestClient(restApiClient *github.Client) *SimpleGitHubRestClient {
	return &SimpleGitHubRestClient{restApiClient: restApiClient}
}

func (simpleGitHubRestClient SimpleGitHubRestClient) ListEvents(ctx context.Context, options *github.ListOptions) ([]*github.Event, error) {
	events, _, err := simpleGitHubRestClient.restApiClient.Activity.ListEvents(ctx, options)
	return events, err
}
