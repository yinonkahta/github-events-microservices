package clients

import (
	"bytes"
	"github-events-microservices/collector/net"
	mocknet "github-events-microservices/collector/net/mocks"
	"github-events-microservices/model"
	"github.com/golang/mock/gomock"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func Test_buildRepoLastUpdatedMap(t *testing.T) {
	type args struct {
		events []model.Event
	}
	now := time.Now()
	tests := []struct {
		name    string
		args    args
		want    map[RepoIdentifier]time.Time
		wantErr bool
	}{
		{
			name: "sanity",
			args: args{[]model.Event{{RepoFullName: "a/b", CreatedAt: now}, {RepoFullName: "c/d", CreatedAt: now}}},
			want: map[RepoIdentifier]time.Time{RepoIdentifier{Owner: "a", Name: "b"}: now, RepoIdentifier{Owner: "c", Name: "d"}: now},
		},
		{
			name: "invalid repo identifier",
			args: args{[]model.Event{{RepoFullName: "ab", CreatedAt: now}, {RepoFullName: "c/d", CreatedAt: now}}},
			want: map[RepoIdentifier]time.Time{RepoIdentifier{Owner: "c", Name: "d"}: now},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildRepoLastUpdatedMap(tt.args.events)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildRepoLastUpdatedMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubGraphQLClient_FetchRepos(t *testing.T) {
	now := time.Now()
	mockCtrl := gomock.NewController(t)
	validResponse := "{\"data\":{\"repoQuery0\":{\"id\":\"R_kgDOJzH7ng\",\"name\":\"RJohnPaul\",\"owner\":{\"login\":\"RJohnPaul\"},\"url\":\"https://github.com/RJohnPaul/RJohnPaul\",\"stargazerCount\":0}}}"

	defer mockCtrl.Finish()
	type fields struct {
		httpClient       net.HttpClient
		githubGraphQLUrl string
		token            string
	}
	type args struct {
		events []model.Event
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Repo
		wantErr bool
	}{
		{
			name: "sanity",
			fields: fields{
				httpClient:       mockHttpClient(mockCtrl, validResponse),
				githubGraphQLUrl: "",
				token:            "",
			},
			args: args{[]model.Event{{RepoFullName: "a/b", CreatedAt: now}, {RepoFullName: "c/d", CreatedAt: now}}},
			want: []model.Repo{{
				ID:            "R_kgDOJzH7ng",
				Owner:         "RJohnPaul",
				Name:          "RJohnPaul",
				Url:           "https://github.com/RJohnPaul/RJohnPaul",
				Stars:         0,
				LastUpdatedAt: time.Time{},
			}},
			wantErr: false,
		},
		{
			name: "empty response",
			fields: fields{
				httpClient:       mockHttpClient(mockCtrl, "{\"data\":{}}"),
				githubGraphQLUrl: "",
				token:            "",
			},
			args:    args{[]model.Event{{RepoFullName: "a/b", CreatedAt: now}, {RepoFullName: "c/d", CreatedAt: now}}},
			want:    []model.Repo{},
			wantErr: false,
		},
		{
			name: "invalid response error",
			fields: fields{
				httpClient:       mockHttpClient(mockCtrl, "hello world"),
				githubGraphQLUrl: "",
				token:            "",
			},
			args:    args{[]model.Event{{RepoFullName: "a/b", CreatedAt: now}, {RepoFullName: "c/d", CreatedAt: now}}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := GithubGraphQLClient{
				httpClient:       tt.fields.httpClient,
				githubGraphQLUrl: tt.fields.githubGraphQLUrl,
				token:            tt.fields.token,
			}
			got, err := receiver.FetchRepos(tt.args.events)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchRepos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchRepos() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockHttpClient(mockCtrl *gomock.Controller, response string) *mocknet.MockHttpClient {
	httpClientMock := mocknet.NewMockHttpClient(mockCtrl)
	httpClientMock.EXPECT().
		Do(gomock.Any()).
		Return(&http.Response{Body: io.NopCloser(bytes.NewBufferString(response))}, nil).
		Times(1)
	return httpClientMock
}
