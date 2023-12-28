package main

import (
	"fmt"
	"github-events-microservices/collector/clients"
	"github-events-microservices/collector/config"
	"github-events-microservices/logging"
	"github-events-microservices/model"
	"github.com/google/go-github/v57/github"
	"log/slog"
	"sync"
	"time"
)

func main() {
	slog.SetDefault(logging.Create())
	slog.Info("Github Events Collector started")

	gitHubClient := newGitHubClient()
	batchStore := clients.NewGithubStoreClient(config.CollectorConfiguration.MongoDbUrl, config.CollectorConfiguration.MongoDbPort)
	eventsChannel := make(chan model.Event)

	var wg sync.WaitGroup
	wg.Add(2)

	go fetchEvents(gitHubClient, eventsChannel, &wg)
	go storeEvents(gitHubClient, batchStore, eventsChannel, &wg)

	wg.Wait()
	close(eventsChannel)
}

func newGitHubClient() *clients.GitHubPublicEventsClient {
	return clients.NewGitHubClient(github.NewClient(nil).WithAuthToken(config.CollectorConfiguration.GitHubToken))
}

func fetchEvents(gitHubClient *clients.GitHubPublicEventsClient, eventsChannel chan model.Event, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		slog.Info("fetching events")
		events, err := fetch(gitHubClient)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to fetch github events: %s", err.Error()))
		}
		for _, event := range events {
			eventsChannel <- event
		}
		time.Sleep(config.CollectorConfiguration.FetchInterval)
	}
}

func fetch(gitHubClient *clients.GitHubPublicEventsClient) ([]model.Event, error) {
	events, err := gitHubClient.ListEvents()
	if err != nil {
		return nil, err
	}
	return events, nil
}

func storeEvents(gitHubClient *clients.GitHubPublicEventsClient, batchStore *clients.GithubStoreClient, eventsChannel chan model.Event, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(config.CollectorConfiguration.MaxTimeout)
	defer ticker.Stop()

	events := make([]model.Event, 0)

	//use a time/size bounded queue to store events/repos/users in batches
	for {
		select {
		case <-ticker.C:
			slog.Debug("reached max timeout")
			if len(events) > 0 {
				slog.Debug(fmt.Sprintf("saving %d items", len(events)))
				repos, err := gitHubClient.FetchRepos(events)
				if err != nil {
					slog.Error(fmt.Sprintf("failed to fetch repos: %s, %v", err.Error(), repos))
				}
				batchStore.Save(events, repos)
				events = make([]model.Event, 0)
			} else {
				slog.Debug("zero items in batch. Skipping saving")
			}
		case event := <-eventsChannel:
			events = append(events, event)
			if len(events) >= config.CollectorConfiguration.MaxItems {
				slog.Debug("reached max items")
				slog.Debug(fmt.Sprintf("saving %d items", len(events)))
				repos, err := gitHubClient.FetchRepos(events)
				if err != nil {
					slog.Error(fmt.Sprintf("failed to fetch repos: %s, %v", err.Error(), repos))
				}
				batchStore.Save(events, repos)
				ticker.Reset(config.CollectorConfiguration.MaxTimeout)
				events = make([]model.Event, 0)
			}
		}
	}
}
