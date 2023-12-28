package config

import (
	"os"
)

var ApiConfiguration *Configuration

type Configuration struct {
	MongoDbUrl       string
	MongoDbPort      string
	EventsDb         string
	EventsCollection string
	ReposDb          string
	ReposCollection  string
	UsersDb          string
	UsersCollection  string
}

func init() {
	ApiConfiguration = &Configuration{
		MongoDbUrl:       getOrDefault(mongoDbUrlKey, defaultMongodbUrl),
		MongoDbPort:      getOrDefault(mongoDbPortKey, defaultMongoDbPort),
		EventsDb:         getOrDefault(eventsDbKey, defaultDb),
		EventsCollection: getOrDefault(eventsCollectionKey, defaultEventsCollection),
		ReposDb:          getOrDefault(reposDbKey, defaultDb),
		ReposCollection:  getOrDefault(reposCollectionKey, defaultReposCollection),
		UsersDb:          getOrDefault(usersDbKey, defaultDb),
		UsersCollection:  getOrDefault(usersCollectionKey, defaultUsersCollection),
	}
}

func getOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
