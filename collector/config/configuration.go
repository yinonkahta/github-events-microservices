package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

var CollectorConfiguration *Configuration

type Configuration struct {
	MongoDbUrl                  string
	MongoDbPort                 string
	GitHubGraphQLRequestTimeout time.Duration
	GitHubToken                 string
	FetchInterval               time.Duration
	MaxItems                    int
	MaxTimeout                  time.Duration
	EventsDb                    string
	EventsCollection            string
	ReposDb                     string
	ReposCollection             string
	UsersDb                     string
	UsersCollection             string
}

func init() {
	CollectorConfiguration = &Configuration{
		MongoDbUrl:                  getOrDefault(mongoDbUrlKey, defaultMongodbUrl),
		MongoDbPort:                 getOrDefault(mongoDbPortKey, defaultMongoDbPort),
		GitHubGraphQLRequestTimeout: getGitHubGraphQLRequestTimeout(),
		GitHubToken:                 getToken(),
		FetchInterval:               getFetchInterval(),
		MaxItems:                    getAsInt(maxItemsKey, defaultMaxItems),
		MaxTimeout:                  getMaxTimeout(),
		EventsDb:                    getOrDefault(eventsDbKey, defaultDb),
		EventsCollection:            getOrDefault(eventsCollectionKey, defaultEventsCollection),
		ReposDb:                     getOrDefault(reposDbKey, defaultDb),
		ReposCollection:             getOrDefault(reposCollectionKey, defaultReposCollection),
		UsersDb:                     getOrDefault(usersDbKey, defaultDb),
		UsersCollection:             getOrDefault(usersCollectionKey, defaultUsersCollection),
	}
}

func getToken() string {
	value := os.Getenv(gitHubToken)
	if len(value) == 0 {
		slog.Error(fmt.Sprintf("Please provide a valid token using th ENV variable: '%s'", gitHubToken))
		os.Exit(1)
	}
	return value
}

func getGitHubGraphQLRequestTimeout() time.Duration {
	fetchInterval := getAsInt(gitHubGraphQLRequestTimeoutSeconds, defaultGraphQLRequestTimeoutSeconds)
	return time.Duration(fetchInterval) * time.Second
}

func getMaxTimeout() time.Duration {
	fetchInterval := getAsInt(maxTimeoutSecondsKey, defaultMaxTimeout)
	return time.Duration(fetchInterval) * time.Second
}

func getFetchInterval() time.Duration {
	fetchInterval := getAsInt(fetchIntervalKey, defaultFetchInterval)
	return time.Duration(fetchInterval) * time.Minute
}

func getOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func getAsInt(key string, defaultValue int) int {
	valueString := os.Getenv(key)
	if len(valueString) == 0 {
		return defaultValue
	}
	value, err := strconv.Atoi(valueString)
	if err != nil {
		slog.Error(fmt.Sprintf("Invalid max items value: %s", valueString))
		os.Exit(1)
	}
	return value
}
