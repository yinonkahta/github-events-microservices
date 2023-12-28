package clients

import (
	"fmt"
	"github-events-microservices/collector/config"
	"github-events-microservices/model"
	"github-events-microservices/stores"
	"log/slog"
	"os"
	"sync"
)

type StoreType string

type GithubStoreClient struct {
	storesMap map[string]stores.ReadWriteStore
	wg        *sync.WaitGroup
}

func (receiver GithubStoreClient) Save(events []model.Event, repos []model.Repo) {
	receiver.wg.Add(3)

	go receiver.saveEvents(events)
	go receiver.saveRepos(repos)
	go receiver.saveUsers(events)

	receiver.wg.Wait()
}

func (receiver GithubStoreClient) saveEvents(events []model.Event) {
	slog.Debug("storing events")
	defer receiver.wg.Done()

	items := make([]interface{}, 0)
	for _, event := range events {
		eventItem := event
		items = append(items, eventItem)
	}

	err := receiver.eventsStore().SaveAll(items)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to save events %s", err.Error()))
	}
}

func (receiver GithubStoreClient) saveRepos(repos []model.Repo) {
	slog.Debug("storing events")
	defer receiver.wg.Done()

	reposMap := make(map[interface{}]interface{})
	for _, repo := range repos {
		reposMap[repo.ID] = repo
	}

	err := receiver.reposStore().UpdateAllById(reposMap)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to save repos %s", err.Error()))
	}
}

func (receiver GithubStoreClient) saveUsers(events []model.Event) {
	slog.Debug("storing users")
	defer receiver.wg.Done()

	usersMap := make(map[interface{}]interface{})
	for _, event := range events {
		user := model.User{
			ID:            event.ActorId,
			Login:         event.ActorLogin,
			Url:           event.ActorUrl,
			AvatarUrl:     event.ActorAvatarUrl,
			LastUpdatedAt: event.CreatedAt,
		}
		usersMap[event.ActorId] = user
	}

	err := receiver.usersStore().UpdateAllById(usersMap)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to save users: %s", err.Error()))
	}
}

func (receiver GithubStoreClient) eventsStore() stores.ReadWriteStore {
	return receiver.storesMap[config.CollectorConfiguration.EventsCollection]
}

func (receiver GithubStoreClient) reposStore() stores.ReadWriteStore {
	return receiver.storesMap[config.CollectorConfiguration.ReposCollection]
}

func (receiver GithubStoreClient) usersStore() stores.ReadWriteStore {
	return receiver.storesMap[config.CollectorConfiguration.UsersCollection]
}

func NewGithubStoreClient(url string, port string) *GithubStoreClient {
	var wg sync.WaitGroup

	storesMap := make(map[string]stores.ReadWriteStore)
	storesMap[config.CollectorConfiguration.EventsCollection] = createMongoStore(url, port, config.CollectorConfiguration.EventsDb, config.CollectorConfiguration.EventsCollection)
	storesMap[config.CollectorConfiguration.ReposCollection] = createMongoStore(url, port, config.CollectorConfiguration.ReposDb, config.CollectorConfiguration.ReposCollection)
	storesMap[config.CollectorConfiguration.UsersCollection] = createMongoStore(url, port, config.CollectorConfiguration.UsersDb, config.CollectorConfiguration.UsersCollection)

	return &GithubStoreClient{
		storesMap: storesMap,
		wg:        &wg,
	}
}

func createMongoStore(url string, port string, database string, collection string) *stores.MongoDbCollectionStore {
	mongoDbStore, err := stores.NewMongoDbStore(url, port, database, collection)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return mongoDbStore
}
