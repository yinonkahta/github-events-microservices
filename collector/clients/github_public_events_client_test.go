package clients

import (
	"encoding/json"
	mockclients "github-events-microservices/collector/clients/mocks"
	"github-events-microservices/model"
	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v57/github"
	"reflect"
	"testing"
	"time"
)

func TestGitHubPublicEventsClient_ListEvents(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type fields struct {
		restApiClient       GitHubRestClient
		githubGraphQLClient GithubGraphQLClient
	}

	id1 := "id1"
	eventType := "type"
	repoName := "repo_name"
	repoUrl := "repo_url"
	userName := "user"
	isPublic := true
	userId := int64(123)
	userUrl := "user_url"
	avatarUrl := "avatar_url"
	tests := []struct {
		name    string
		fields  fields
		want    []model.Event
		wantErr bool
	}{
		{
			name: "non empty results",
			fields: fields{
				restApiClient: mockRestApiClient(mockCtrl, []*github.Event{{
					Type:       &eventType,
					Public:     &isPublic,
					RawPayload: &json.RawMessage{},
					Repo:       &github.Repository{Name: &repoName, URL: &repoUrl},
					Actor:      &github.User{ID: &userId, Login: &userName, URL: &userUrl, AvatarURL: &avatarUrl},
					Org:        &github.Organization{},
					CreatedAt:  &github.Timestamp{},
					ID:         &id1,
				}}, nil),
				githubGraphQLClient: GithubGraphQLClient{},
			},
			want: []model.Event{{
				ID:             id1,
				Type:           eventType,
				CreatedAt:      time.Time{},
				Public:         true,
				RepoFullName:   repoName,
				RepoUrl:        repoUrl,
				ActorLogin:     userName,
				ActorId:        userId,
				ActorUrl:       userUrl,
				ActorAvatarUrl: avatarUrl,
			}},
			wantErr: false,
		},
		{
			name: "empty results",
			fields: fields{
				restApiClient:       mockRestApiClient(mockCtrl, []*github.Event{}, nil),
				githubGraphQLClient: GithubGraphQLClient{},
			},
			want:    []model.Event{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := GitHubPublicEventsClient{
				restApiClient:       tt.fields.restApiClient,
				githubGraphQLClient: tt.fields.githubGraphQLClient,
			}
			got, err := receiver.ListEvents()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListEvents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListEvents() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockRestApiClient(mockCtrl *gomock.Controller, events []*github.Event, err error) GitHubRestClient {
	githubRestClientMock := mockclients.NewMockGitHubRestClient(mockCtrl)
	githubRestClientMock.EXPECT().
		ListEvents(gomock.Any(), gomock.Any()).
		Return(events, err).
		Times(1)
	return githubRestClientMock
}
