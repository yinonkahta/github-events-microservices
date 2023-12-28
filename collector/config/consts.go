package config

const (
	// env variables
	mongoDbUrlKey                      = "MONGO_DB_URL"
	mongoDbPortKey                     = "MONGO_DB_PORT"
	gitHubGraphQLRequestTimeoutSeconds = "GITHUB_GRAPHQL_REQUEST_TIMEOUT"
	gitHubToken                        = "GITHUB_TOKEN"
	fetchIntervalKey                   = "FETCH_INTERVAL_MINUTES"
	maxItemsKey                        = "MAX_ITEMS"
	maxTimeoutSecondsKey               = "MAX_TIMEOUT_SECONDS"
	eventsDbKey                        = "EVENTS_DB"
	eventsCollectionKey                = "EVENTS_COLLECTION"
	reposDbKey                         = "REPOS_DB"
	reposCollectionKey                 = "REPOS_COLLECTION"
	usersDbKey                         = "USERS_DB"
	usersCollectionKey                 = "USERS_COLLECTION"

	defaultMongodbUrl                   = "localhost"
	defaultMongoDbPort                  = "27017"
	defaultGraphQLRequestTimeoutSeconds = 30
	defaultFetchInterval                = 1
	defaultMaxItems                     = 3
	defaultMaxTimeout                   = 10

	defaultDb               = "github"
	defaultEventsCollection = "events"
	defaultReposCollection  = "repos"
	defaultUsersCollection  = "users"
)
