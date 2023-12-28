package config

const (
	// env variables
	mongoDbUrlKey       = "MONGO_DB_URL"
	mongoDbPortKey      = "MONGO_DB_PORT"
	eventsDbKey         = "EVENTS_DB"
	eventsCollectionKey = "EVENTS_COLLECTION"
	reposDbKey          = "REPOS_DB"
	reposCollectionKey  = "REPOS_COLLECTION"
	usersDbKey          = "USERS_DB"
	usersCollectionKey  = "USERS_COLLECTION"

	defaultMongodbUrl  = "localhost"
	defaultMongoDbPort = "27017"

	defaultDb               = "github"
	defaultEventsCollection = "events"
	defaultReposCollection  = "repos"
	defaultUsersCollection  = "users"
	DataType                = "dataType"
	DefaultDataType         = "events"
	LimitParamKey           = "limit"
	DefaultLimit            = 20
	OrderByColumnQueryParam = "orderBy"
	OrderTypeQueryParam     = "orderType"
	DefaultOrderByColumn    = "_id"
	Ascending               = "ascending"
	Descending              = "descending"
)
