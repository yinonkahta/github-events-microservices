package stores

// StubStore for tests
type StubStore struct {
	data []interface{}
}

func (store *StubStore) All(results interface{}) error {
	return store.Get(0, OrderBy{
		Column: "_id",
		Order:  1,
	}, results)
}

func (store *StubStore) Get(limit int64, orderBy OrderBy, results interface{}) error {
	resultsArr := results.([]interface{})
	for i := 0; i < len(resultsArr); i++ {
		if len(store.data) == i {
			break
		}
		resultsArr[i] = store.data[i]
	}
	return nil
}

func (store *StubStore) Count() (int64, error) {
	return int64(len(store.data)), nil
}

func (store *StubStore) Close() error {
	store.data = nil
	return nil
}

func (store *StubStore) Save(element interface{}) error {
	store.data = append(store.data, element)
	return nil
}

func (store *StubStore) SaveAll(elements []interface{}) error {
	for _, element := range elements {
		store.data = append(store.data, element)
	}
	return nil
}

func (store *StubStore) UpdateAllById(elements map[interface{}]interface{}) error {
	for _, element := range elements {
		store.data = append(store.data, element)
	}
	return nil
}

func NewStubStore(data []interface{}) *StubStore {
	return &StubStore{data: data}
}
