package clients

import (
	"github-events-microservices/collector/config"
	"github-events-microservices/model"
	"github-events-microservices/stores"
	"reflect"
	"sync"
	"testing"
)

func TestGithubStoreClient_Save(t *testing.T) {
	type fields struct {
		storesMap map[string]stores.ReadWriteStore
		wg        *sync.WaitGroup
	}
	type args struct {
		events []model.Event
		repos  []model.Repo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "sanity",
			fields: fields{
				storesMap: createStoresMap(),
				wg:        &sync.WaitGroup{},
			},
			args: args{
				events: []model.Event{{ID: "ev1", ActorId: 1}, {ID: "ev2", ActorId: 2}},
				repos:  []model.Repo{{ID: "rp1"}, {ID: "rp2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := GithubStoreClient{
				storesMap: tt.fields.storesMap,
				wg:        tt.fields.wg,
			}
			receiver.Save(tt.args.events, tt.args.repos)
			events := make([]interface{}, 2)
			events[0] = model.Event{ID: "ev1", ActorId: 1}
			events[1] = model.Event{ID: "ev2", ActorId: 2}

			wantedEvents := stores.NewStubStore(events)
			if !reflect.DeepEqual(receiver.eventsStore(), wantedEvents) {
				t.Errorf("eventsStore() = %v, want %v", receiver.eventsStore(), wantedEvents)
			}

			repos := make([]interface{}, 2)
			repos[0] = model.Repo{ID: "rp1"}
			repos[1] = model.Repo{ID: "rp2"}

			wantedRepos := stores.NewStubStore(repos)
			if !reflect.DeepEqual(receiver.reposStore(), wantedRepos) {
				t.Errorf("reposStore() = %v, want %v", receiver.reposStore(), wantedRepos)
			}

			users := make([]interface{}, 2)
			users[0] = model.User{ID: 1}
			users[1] = model.User{ID: 2}

			wantedUsers := stores.NewStubStore(users)
			if !reflect.DeepEqual(receiver.usersStore(), wantedUsers) {
				t.Errorf("usersStore() = %v, want %v", receiver.usersStore(), wantedUsers)
			}
		})
	}
}

func createStoresMap() map[string]stores.ReadWriteStore {
	storesMap := make(map[string]stores.ReadWriteStore)
	var eventsData []interface{}
	var reposData []interface{}
	var usersData []interface{}
	storesMap[config.CollectorConfiguration.EventsCollection] = stores.NewStubStore(eventsData)
	storesMap[config.CollectorConfiguration.ReposCollection] = stores.NewStubStore(reposData)
	storesMap[config.CollectorConfiguration.UsersCollection] = stores.NewStubStore(usersData)
	return storesMap
}
